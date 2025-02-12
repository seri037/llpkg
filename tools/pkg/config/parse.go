package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

func ParseLLpkgConfig(configPath string) (LLpkgConfig, error) {
	var config LLpkgConfig
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
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
		return config, errors.New("invalid configuration: package.name is required")
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
