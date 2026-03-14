package main

import "runtime"

const (
	defaultMinVersion = "1.0.0"

	defaultGoBuildOutputFolder = ".gobuild"
	defaultGoBuildConfig       = ".gobuild-config"
	defaultGoBuildVersionFile  = ".gobuild-version"
	defaultGoBuildBinaryFile   = ".gobuild-binary"
)

var (
	version = "1.0.1"

	// go tool dist list
	// https://go.dev/wiki/GoArm
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
	Compress        bool     `json:"compress,omitempty"`
	Compressor      string   `json:"compressor,omitempty"`
	CompressorFlags []string `json:"compressorFlags,omitempty"`
	Platforms       []string `json:"platforms,omitempty"`
}

func NewConfig() Config {
	config := Config{
		Compress:   true,
		Compressor: "upx",
		Platforms:  defaultPlatforms,
	}

	if runtime.GOOS == "windows" {
		config.Compressor += ".exe"
	}

	return config
}

type BuildConfig struct {
	Binary      string
	Config      Config
	Version     string
	GoModFolder string
}
