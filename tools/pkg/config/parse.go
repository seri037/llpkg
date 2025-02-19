package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime/debug"
	"slices"
	"strings"
)

type ConanSearchResult struct {
	Conancenter map[string]struct{} `json:"conancenter"`
}

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

	spinner := NewLoadingSpinner("Searching for available versions")
	spinner.Start()

	cmd := exec.Command("conan", "search", config.Package.Name, "-r", "conancenter")
	out, err := cmd.Output()

	spinner.Stop()

	if err != nil {
		return config, withStack(fmt.Errorf("failed to execute conan command: %w", err))
	}
	cmdString := string(out)
	fmt.Print(cmdString)
	versions := extractVersions(cmdString, config.Package.Name)
	if len(versions) == 0 {
		// fallback to json output
		var result ConanSearchResult
		cmd := exec.Command("conan", "search", config.Package.Name, "-r", "conancenter", "-f", "json")
		out, err := cmd.Output()
		if err != nil {
			return config, withStack(fmt.Errorf("failed to execute conan command: %w", err))
		}
		err = json.Unmarshal(out, &result)
		if err != nil {
			return config, withStack(fmt.Errorf("failed to decode json output: %w", err))
		}
		for versionString := range result.Conancenter {
			versions = append(versions, strings.Split(versionString, "/")[1])
		}
	}
	if len(versions) == 0 {
		return config, withStack(errors.New("no versions found"))
	}

	// check if the package name is set
	if config.Package.Name == "" {
		return config, withStack(errors.New("invalid configuration: package.name is required"))
	}

	if config.Package.CVersion == "" {
		fmt.Println("Warning: package.cVersion is not set ")
		config.Package.CVersion = selectVersion(versions)
		if config.Package.CVersion == "" {
			return config, withStack(errors.New("invalid version"))
		}
	} else {
		if !slices.Contains(versions, config.Package.CVersion) {
			fmt.Print("Your input version is not in the list")
			fmt.Println("Available versions:", strings.Join(versions, ", "))
			return config, withStack(errors.New("invalid version"))
		}
	}

	var moduleVersionRegex = regexp.MustCompile(`v\d+\.\d+\.\d+`)
	if config.Package.ModuleVersion == "" {
		fmt.Println("Warning: package.moduleVersion is not set.")
		var moduleVersion string
		for !moduleVersionRegex.MatchString(moduleVersion) {
			fmt.Println("Input the SEMANTIC version with \"v\" you would like to use (eg: v1.0.0): ")
			reader := bufio.NewReader(os.Stdin)
			moduleVersion, _ = reader.ReadString('\n')
			moduleVersion = strings.TrimSpace(moduleVersion)
		}
	} else {
		if !moduleVersionRegex.MatchString(config.Package.ModuleVersion) {
			fmt.Println("Invalid module version")
			fmt.Println("Filed \"moduleVersion\" requires a SEMANTIC version with \"v\" (eg: v1.0.0)")
			return config, withStack(errors.New("invalid module version"))
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

func extractVersions(consoleOutput string, pkgName string) []string {
	// 按行分割控制台输出
	lines := strings.Split(consoleOutput, "\n")
	pkgNameRegex := regexp.MustCompile(pkgName + "/" + ".*")

	// 定义一个切片来存储版本号
	var versions []string

	// 从最后一行开始反向遍历
	for i := len(lines) - 1; i >= 0; i-- {
		// 匹配包名
		if pkgNameRegex.MatchString(lines[i]) {
			versions = append([]string{strings.Split(lines[i], "/")[1]}, versions...)
		} else {
			break
		}
	}

	return versions
}

func selectVersion(versions []string) string {
	fmt.Println("Available versions:", strings.Join(versions, ", "))
	fmt.Print("Which version would you like to use? (n to exit, RETURN to use latest): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	input = strings.TrimSpace(input)
	if input == "n" || input == "N" || input == "no" || input == "No" {
		return ""
	}

	if input == "" {
		return versions[len(versions)-1]
	} else if !slices.Contains(versions, input) {
		return ""
	} else {
		return input
	}
}
