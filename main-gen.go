package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func gen() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("cwd: %s\n", cwd)

	// .gobuild
	var config = Config{
		Platforms: defaultPlatforms,
	}

	configAsBytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	goBuildFile := filepath.Join(cwd, defaultGoBuildFile)

	err = os.WriteFile(goBuildFile, configAsBytes, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("gen: %s\n", goBuildFile)

	// .gobuild-version
	goBuildVersionFile := filepath.Join(cwd, defaultGoBuildVersionFile)

	err = os.WriteFile(goBuildVersionFile, []byte(defaultMinVersion), os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("gen: %s\n", goBuildVersionFile)

	// .gobuild-binary
	goBuildBinaryFile := filepath.Join(cwd, defaultGoBuildBinaryFile)

	err = os.WriteFile(goBuildBinaryFile, []byte(filepath.Base(cwd)), os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Printf("gen: %s\n", goBuildBinaryFile)
}
