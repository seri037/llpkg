package config

import (
	"testing"
)

func TestParseLLpkgConfig(t *testing.T) {
	config, err := ParseLLpkgConfig("../../demo/.llpkg/llpkglite.cfg")
	if err != nil {
		t.Errorf("Error parsing config file: %v", err)
	}
	PrintStruct(config, "")
}
