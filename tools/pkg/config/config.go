package config

type LLpkgConfig struct {
	Upstream Upstream `json:"upstream"`
}

type Upstream struct {
	Installer Installer `json:"installer"`
	Package   Package   `json:"package"`
}

type Installer struct {
	Name   string `json:"name"`
	Config Config `json:"config,omitempty"`
}

type Config struct {
	Options string `json:"options,omitempty"`
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
