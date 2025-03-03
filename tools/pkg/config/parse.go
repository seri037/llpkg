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
	if config.Upstream.Installer.Name == "" {
		config.Upstream.Installer.Name = "conan"
	}
	if config.Upstream.Installer.Config.Options == "" {
		config.Upstream.Installer.Config.Options = ""
	}

	spinner := NewLoadingSpinner("Searching for available versions")
	spinner.Start()

	cmd := exec.Command("conan", "search", config.Upstream.Package.Name, "-r", "conancenter")
	out, err := cmd.Output()

	spinner.Stop()

	if err != nil {
		return config, withStack(fmt.Errorf("failed to execute conan command: %w", err))
	}
	cmdString := string(out)
	fmt.Print(cmdString)
	versions := extractVersions(cmdString, config.Upstream.Package.Name)
	if len(versions) == 0 {
		// fallback to json output
		var result ConanSearchResult
		cmd := exec.Command("conan", "search", config.Upstream.Package.Name, "-r", "conancenter", "-f", "json")
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
	if config.Upstream.Package.Name == "" {
		return config, withStack(errors.New("invalid configuration: package.name is required"))
	}

	if config.Upstream.Package.Version == "" {
		fmt.Println("Warning: package.cVersion is not set ")
		config.Upstream.Package.Version = selectVersion(versions)
		if config.Upstream.Package.Version == "" {
			return config, withStack(errors.New("invalid version"))
		}
	} else {
		if !slices.Contains(versions, config.Upstream.Package.Version) {
			fmt.Println("Your input version is not in the list")
			fmt.Println("Available versions:", strings.Join(versions, ", "))
			return config, withStack(errors.New("invalid version"))
		}
	}

	// var moduleVersionRegex = regexp.MustCompile(`v\d+\.\d+\.\d+`)
	// if config.Package.ModuleVersion == "" {
	// 	fmt.Println("Warning: package.moduleVersion is not set.")
	// 	var moduleVersion string
	// 	for !moduleVersionRegex.MatchString(moduleVersion) {
	// 		fmt.Println("Input the SEMANTIC version with \"v\" you would like to use (eg: v1.0.0): ")
	// 		reader := bufio.NewReader(os.Stdin)
	// 		moduleVersion, _ = reader.ReadString('\n')
	// 		moduleVersion = strings.TrimSpace(moduleVersion)
	// 	}
	// } else {
	// 	if !moduleVersionRegex.MatchString(config.Package.ModuleVersion) {
	// 		fmt.Println("Invalid module version")
	// 		fmt.Println("Filed \"moduleVersion\" requires a SEMANTIC version with \"v\" (eg: v1.0.0)")
	// 		return config, withStack(errors.New("invalid module version"))
	// 	}
	// }

	return config, nil
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
	lines := strings.Split(consoleOutput, "\n")
	versionPattern := regexp.MustCompile(`\s+` + regexp.QuoteMeta(pkgName) + `.*`)

	var versions []string
	var inPackageSection bool
	var currentIndent int

	for _, line := range lines {
		trimmedLine := strings.TrimLeft(line, " ")
		indent := len(line) - len(trimmedLine)

		if strings.TrimSpace(line) == pkgName {
			inPackageSection = true
			currentIndent = indent
			continue
		}

		if inPackageSection && indent <= currentIndent {
			inPackageSection = false
		}

		if inPackageSection && indent > currentIndent {
			matches := versionPattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				versions = append(versions, matches[1])
			}
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
