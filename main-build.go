package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var (
	goConventionsFolderLayouts = []string{
		"pkg",
	}
)

func build() {
	cwd, err := os.Getwd()
	exitIfErr(err)

	productVersion := "1.0.0"
	productVersionRaw, err := os.ReadFile(defaultGoBuildVersionFile)
	if err == nil {
		productVersion = strings.TrimSpace(string(productVersionRaw))
	}

	productBinaryName := ""
	productBinaryNameRaw, err := os.ReadFile(defaultGoBuildBinaryFile)
	if err != nil {
		productBinaryName = filepath.Base(cwd)                              // if file name not set then use folder name
		if slices.Contains(goConventionsFolderLayouts, productBinaryName) { // account for cases when code is in "pkg" type folders
			productBinaryName = filepath.Base(filepath.Dir(cwd))
		}
	} else {
		productBinaryName = strings.TrimSpace(string(productBinaryNameRaw))
	}
	productBinaryName += "_v" + productVersion

	productBuildConfig := Config{
		Platforms: defaultPlatforms,
	}
	productBuildConfigRaw, err := os.ReadFile(defaultGoBuildFile)
	if err == nil {
		err = json.Unmarshal(productBuildConfigRaw, &productBuildConfig)
		exitIfErr(err)
	}

	productOutputFolder := filepath.Join(cwd, defaultGoBuildOutputFolder)

	log.Println()
	log.Printf("build  : %s @ %s\n", productBinaryName, productVersion)
	log.Printf("output : %s\n", productOutputFolder)
	log.Printf("config : %s\n", productBuildConfig.Platforms)

	err = os.RemoveAll(productOutputFolder)
	exitIfErr(err)

	err = os.MkdirAll(productOutputFolder, os.FileMode(0755))
	exitIfErr(err)

	var debugFlags = ""
	if false {
		debugFlags = "-tags debug"
	}

	var ldFlags = fmt.Sprintf("-w -s -X main.version=\"%s\"", productVersion)

	// format: $binary.$arch-$details-$os.$ext
	for _, platform := range productBuildConfig.Platforms {
		// get OS + ARCH
		productOSAndArch := strings.Split(platform, "/")

		productOS := ""
		productOSLabel := ""
		if len(productOSAndArch) > 0 {
			productOS = productOSAndArch[0]
			productOSLabel = productOS
			if _, hasKey := osToDisplay[productOS]; hasKey {
				productOSLabel = osToDisplay[productOS]
			}
		}

		productArch := ""
		productArchLabel := ""
		if len(productOSAndArch) > 1 {
			productArch = productOSAndArch[1]
			productArchLabel = productArch
			if _, hasKey := archToDisplay[productArch]; hasKey {
				productArchLabel = archToDisplay[productArch]
			}
		}

		productArmVersion := ""
		if len(productOSAndArch) > 2 {
			productArmVersion = productOSAndArch[2]
		}

		productArchWithDetails := productArchLabel
		if productArmVersion != "" {
			productArchWithDetails = fmt.Sprintf("%s-v%s", productArchLabel, productArmVersion)
		}

		productExt := ""
		if _, hasKey := osToExt[productOS]; hasKey {
			productExt = osToExt[productOS]
		}

		outputFile := filepath.Join(productOutputFolder, fmt.Sprintf("%s_%s-%s%s",
			productBinaryName,
			productArchWithDetails,
			productOSLabel,
			productExt,
		))

		log.Printf("\t-> %-55s \n", filepath.Base(outputFile))

		cmdToRunEnv := os.Environ() // copy current env

		for productEnvK, productEnvV := range defaultEnvs {
			cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("%s=%s", productEnvK, productEnvV))
		}

		cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("GOOS=%s", productOS))
		cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("GOARCH=%s", productArch))
		if productArmVersion != "" {
			cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("GOARM=%s", productArmVersion))
		}

		cmdToRun := exec.Command("go", "build")
		if debugFlags != "" {
			cmdToRun.Args = append(cmdToRun.Args, debugFlags)
		}
		cmdToRun.Args = append(cmdToRun.Args, "-ldflags", ldFlags)
		cmdToRun.Args = append(cmdToRun.Args, "-o", outputFile)
		cmdToRun.Args = append(cmdToRun.Args, "./")
		cmdToRun.Env = cmdToRunEnv
		cmdToRun.Dir = cwd
		cmdToRun.Stderr = os.Stderr
		cmdToRun.Stderr = os.Stdout
		err = cmdToRun.Run()
		exitIfErr(err)
	}
}
