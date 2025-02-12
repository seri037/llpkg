package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime/debug"
	"slices"
	"strings"
)

func withStack(err error) error {
	return fmt.Errorf("%w\n%s", err, debug.Stack())
}

func ParseLLpkgConfig(configPath string) (LLpkgConfig, error) {
	var config LLpkgConfig
	file, err := os.Open(configPath)
	if err != nil {
		return config, withStack(fmt.Errorf("failed to open config file: %w", err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, withStack(fmt.Errorf("failed to decode config file: %w", err))
	}

	// set default values
	if !config.Package.VersionChange {
		config.Package.VersionChange = true
	}
	if config.Upstream.Name == "" {
		config.Upstream.Name = "conan"
	}
	if config.Upstream.Config.Options == "" {
		config.Upstream.Config.Options = ""
	}
	if config.Toolchain.Name == "" {
		config.Toolchain.Name = "llcppg"
	}
	if config.Toolchain.Version == "" {
		config.Toolchain.Version = "latest"
	}

	// check if the package name is set
	if config.Package.Name == "" {
		return config, withStack(errors.New("invalid configuration: package.name is required"))
	}
	if config.Package.Version == "" {
		fmt.Println("Warning: package.version is not set \n Setting it to latest ")

		cmd := exec.Command("conan", "search", config.Package.Name, "-f", "json", "-r", "conancenter")
		out, err := cmd.Output()
		if err != nil {
			return config, withStack(fmt.Errorf("failed to execute conan command: %w", err))
		}

		scanner := bufio.NewScanner(bytes.NewReader(out))
		var latestVersion string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, `"`+config.Package.Name+`/`) {
				line = strings.Trim(line, `",`)
				parts := strings.Split(line, "/")
				if len(parts) != 2 {
					continue
				}
				latestVersion = parts[1]
			}
		}

		versionInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(out), &versionInfo)
		if err != nil {
			return config, withStack(fmt.Errorf("failed to parse JSON response: %w", err))
		}

		availableVersion := make([]string, 0, len(versionInfo["conancenter"].(map[string]interface{})))
		for k := range versionInfo["conancenter"].(map[string]interface{}) {
			availableVersion = append(availableVersion, strings.Split(k, "/")[1])
		}
		fmt.Println("Available versions:", strings.Join(availableVersion, ", "))
		fmt.Print("Which version would you like to use? (n to exit, return to use latest): ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return config, withStack(fmt.Errorf("failed to read user input: %w", err))
		}
		input = strings.TrimSpace(input)
		if input == "n" || input == "N" || input == "no" || input == "No" {
			return config, withStack(errors.New("package.version is required"))
		}

		if input == "" {
			config.Package.Version = latestVersion
		} else if !slices.Contains(availableVersion, input) {
			return config, withStack(errors.New("invalid version"))
		} else {
			config.Package.Version = input
		}
	}

	// check if the config is valid
	err = validateConfig(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func validateConfig(config LLpkgConfig) error {
	v := reflect.ValueOf(config)
	t := reflect.TypeOf(config)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			for j := 0; j < field.NumField(); j++ {
				subField := field.Field(j)
				subFieldType := fieldType.Type.Field(j)

				if !subField.IsValid() {
					return fmt.Errorf("invalid configuration: %s.%s is required", fieldType.Name, subFieldType.Name)
				}
			}
		} else {
			if !field.IsValid() {
				return fmt.Errorf("invalid configuration: %s is required", fieldType.Name)
			}
		}
	}

	return nil
}

func PrintStruct(s interface{}, indent string) {
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		if fieldVal.Kind() == reflect.Struct {
			fmt.Printf("%s%s:\n", indent, fieldType.Name)
			PrintStruct(fieldVal.Interface(), indent+"  ")
		} else {
			fmt.Printf("%s%s: %v\n", indent, fieldType.Name, fieldVal.Interface())
		}
	}
}
