package main

const (
	defaultGoBuildFile         = ".gobuild"
	defaultGoBuildVersionFile  = ".gobuild-version"
	defaultGoBuildBinaryFile   = ".gobuild-binary"
	defaultGoBuildOutputFolder = ".gobuild-"

	defaultMinVersion = "1.0.0"
)

var (
	version = "1.0.0"

	defaultPlatforms = []string{
		"linux/amd64",
		"linux/arm/5",
		"linux/arm/6",
		"linux/arm/7",
		"linux/arm64",
		"darwin/amd64",
		"darwin/arm64",
		"windows/386",
		"windows/amd64",
		"windows/arm64",
	}

	defaultEnvs = map[string]string{
		"CGO_ENABLED": "0",
		"GO111MODULE": "on",
	}

	osToDisplay = map[string]string{
		"darwin": "osx",
	}

	osToExt = map[string]string{
		"windows": ".exe",
	}

	archToDisplay = map[string]string{
		"386":   "x86",
		"amd64": "x86_64",
		"arm64": "aarch64",
	}
)

type Config struct {
	Platforms []string `json:"platforms,omitempty"`
}
