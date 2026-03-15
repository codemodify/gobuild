package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	goConventionsFolderLayouts = []string{
		"pkg",
	}
)

func gen() {
	buildConfig, cwd := genBuildConfig()

	// binary
	goBuildBinaryFile := filepath.Join(cwd, defaultGoBuildBinaryFile)

	err := os.MkdirAll(filepath.Dir(goBuildBinaryFile), os.FileMode(0755))
	exitIfErr(err)

	err = os.WriteFile(goBuildBinaryFile, []byte(buildConfig.Binary), os.FileMode(0644))
	exitIfErr(err)

	log.Printf("gen: %s\n", goBuildBinaryFile)

	// config
	configAsBytes, err := json.MarshalIndent(buildConfig.Config, "", "\t")
	exitIfErr(err)

	goBuildConfigFile := filepath.Join(cwd, defaultGoBuildConfig)

	err = os.MkdirAll(filepath.Dir(goBuildConfigFile), os.FileMode(0755))
	exitIfErr(err)

	err = os.WriteFile(goBuildConfigFile, configAsBytes, os.FileMode(0644))
	exitIfErr(err)

	log.Printf("gen: %s\n", goBuildConfigFile)

	// version
	goBuildVersionFile := filepath.Join(cwd, defaultGoBuildVersionFile)

	err = os.MkdirAll(filepath.Dir(goBuildVersionFile), os.FileMode(0755))
	exitIfErr(err)

	err = os.WriteFile(goBuildVersionFile, []byte(buildConfig.Version), os.FileMode(0644))
	exitIfErr(err)

	log.Printf("gen: %s\n", goBuildVersionFile)
}

func genRM() {
	cwd, err := os.Getwd()
	exitIfErr(err)

	// binary
	goBuildBinaryFile := filepath.Join(cwd, defaultGoBuildBinaryFile)
	os.Remove(goBuildBinaryFile)

	// config
	goBuildConfigFile := filepath.Join(cwd, defaultGoBuildConfig)
	os.Remove(goBuildConfigFile)

	// version
	goBuildVersionFile := filepath.Join(cwd, defaultGoBuildVersionFile)
	os.Remove(goBuildVersionFile)
}

func genBuildConfig() (BuildConfig, string) {
	cwd, err := os.Getwd()
	exitIfErr(err)

	goModFolder, err := findGoModFolder(cwd)
	exitIfErr(err)

	// binary
	binary := filepath.Base(goModFolder)
	if slices.Contains(goConventionsFolderLayouts, binary) { // account for cases when code is in "pkg" type folders
		binary = filepath.Base(filepath.Dir(goModFolder))
	}

	if raw, err := os.ReadFile(defaultGoBuildBinaryFile); err == nil {
		binary = strings.TrimSpace(string(raw))
	}

	// config
	config := NewConfig()
	if raw, err := os.ReadFile(defaultGoBuildConfig); err == nil {
		err = json.Unmarshal(raw, &config)
		exitIfErr(err)
	}

	// version
	version := defaultMinVersion
	if raw, err := os.ReadFile(defaultGoBuildVersionFile); err == nil {
		version = strings.TrimSpace(string(raw))
	}

	return BuildConfig{
		Binary:      binary,
		Config:      config,
		Version:     version,
		GoModFolder: goModFolder,
	}, cwd
}

func findGoModFolder(cwd string) (string, error) {
	foldersToCheck := []string{
		cwd,
		filepath.Join(cwd, "pkg"),
		filepath.Join(cwd, "src"),
	}

	for _, folder := range foldersToCheck {
		if isFile(filepath.Join(folder, "go.mod")) {
			return folder, nil
		}
	}

	return "", fmt.Errorf("go.mod folder not found")
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return !info.IsDir()
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}
