package config

type LLpkgConfig struct {
	Upstream  Upstream  `json:"upstream"`
	Generator Generator `json:"generator"`
}

type Upstream struct {
	Installer string         `json:"installer"`
	Config    UpstreamConfig `json:"config"`
	Package   Package        `json:"package"`
}

type UpstreamConfig struct {
	Options string `json:"options"`
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Generator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
