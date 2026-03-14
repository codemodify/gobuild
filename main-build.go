package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func build() {
	buildConfig, cwd := genBuildConfig()

	productOutputFolder := filepath.Join(cwd, defaultGoBuildOutputFolder)

	log.Println()
	log.Printf("output   : %s\n", productOutputFolder)
	log.Printf("config   : %v\n", buildConfig.Config)
	log.Printf("build    : %s @ %s\n", buildConfig.Binary, buildConfig.Version)

	err := os.RemoveAll(productOutputFolder)
	exitIfErr(err)

	err = os.MkdirAll(productOutputFolder, os.FileMode(0755))
	exitIfErr(err)

	var debugFlags = ""
	if false {
		debugFlags = "-tags debug"
	}

	var ldFlags = fmt.Sprintf("-s -w -X main.version=\"%s\"", buildConfig.Version)

	// format: $binary_$version-$os.$arch-$details.$ext
	for _, platform := range buildConfig.Config.Platforms {
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

		productBinaryName := buildConfig.Binary + "_v" + buildConfig.Version
		outputFile := filepath.Join(productOutputFolder, fmt.Sprintf("%s_%s-%s%s",
			productBinaryName,
			productOSLabel,
			productArchWithDetails,
			productExt,
		))

		log.Printf("        -> %-55s \n", filepath.Base(outputFile))

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
		cmdToRun.Dir = buildConfig.GoModFolder
		cmdToRun.Stderr = os.Stderr
		cmdToRun.Stdout = os.Stdout
		err = cmdToRun.Run()
		exitIfErr(err)
	}

	// run compressor
	if buildConfig.Config.Compress {
		if _, err := exec.LookPath(buildConfig.Config.Compressor); err == nil {
			filepath.WalkDir(productOutputFolder, func(path string, d fs.DirEntry, err error) error {
				if !isFile(path) {
					return nil
				}

				preSize := getFileSize(path)
				preSizeKB := preSize / 1024
				preSizeMB := preSize / (1024 * 1024)

				cmdToRun := exec.Command(buildConfig.Config.Compressor)
				if len(buildConfig.Config.CompressorFlags) > 0 {
					cmdToRun.Args = append(cmdToRun.Args, buildConfig.Config.CompressorFlags...)
				}
				cmdToRun.Args = append(cmdToRun.Args, path)

				_ = cmdToRun.Run()

				postSize := getFileSize(path)
				postSizeKB := postSize / 1024
				postSizeMB := postSize / (1024 * 1024)

				log.Printf("compress : %s %d/%dK/%dMB -> %d/%dK/%dMB \n", d.Name(), preSize, preSizeKB, preSizeMB, postSize, postSizeKB, postSizeMB)

				return nil
			})
		}
	}
}
