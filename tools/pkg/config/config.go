package config

type LLpkgConfig struct {
	Package   Package   `json:"package"`
	Upstream  Upstream  `json:"upstream,omitempty"`
	Toolchain Toolchain `json:"toolchain,omitempty"`
}

type Package struct {
	Name          string `json:"name"`
	CVersion      string `json:"cVersion,omitempty"`
	ModuleVersion string `json:"moduleVersion,omitempty"`
}

type Upstream struct {
	Name   string         `json:"name,omitempty"`
	Config UpstreamConfig `json:"config,omitempty"`
}

type UpstreamConfig struct {
	Options string `json:"options,omitempty"`
}

type Toolchain struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}
