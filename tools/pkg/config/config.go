package config

type LLpkgConfig struct {
	Package   Package   `json:"package"`
	Upstream  Upstream  `json:"upstream"`
	Toolchain Toolchain `json:"toolchain"`
}

type Package struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	VersionChange bool   `json:"versionChange"`
}

type Upstream struct {
	Name   string         `json:"name"`
	Config UpstreamConfig `json:"config"`
}

type UpstreamConfig struct {
	Options string `json:"options"`
}

type Toolchain struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
