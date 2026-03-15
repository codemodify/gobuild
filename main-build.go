package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildClean() {
	cwd, err := os.Getwd()
	exitIfErr(err)

	productOutputFolder := filepath.Join(cwd, defaultGoBuildOutputFolder)

	err = os.RemoveAll(productOutputFolder)
	exitIfErr(err)
}

func build() {
	buildClean()

	buildConfig, cwd := genBuildConfig()

	productOutputFolder := filepath.Join(cwd, defaultGoBuildOutputFolder)

	log.Println()
	log.Printf("output   : %s\n", productOutputFolder)
	log.Printf("config   : %v\n", buildConfig.Config)
	log.Printf("build    : %s @ %s\n", buildConfig.Binary, buildConfig.Version)

	err := os.MkdirAll(productOutputFolder, os.FileMode(0755))
	exitIfErr(err)

	var ldFlags = fmt.Sprintf("-s -w -X main.version=\"%s\"", buildConfig.Version)

	// format: $binary_$version-$os.$arch-$details.$ext
	for _, platform := range buildConfig.Config.Platforms {
		target, err := parsePlatform(platform)
		exitIfErr(err)

		productBinaryName := buildConfig.Binary + "_v" + buildConfig.Version
		outputFile := filepath.Join(productOutputFolder, fmt.Sprintf("%s_%s-%s%s",
			productBinaryName,
			target.productOSLabel,
			target.productArchWithDetails,
			target.productExt,
		))

		log.Printf("        -> %-55s \n", filepath.Base(outputFile))

		cmdToRunEnv := os.Environ() // copy current env

		for productEnvK, productEnvV := range defaultEnvs {
			cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("%s=%s", productEnvK, productEnvV))
		}

		cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("GOOS=%s", target.productOS))
		cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("GOARCH=%s", target.productArch))
		if target.productArmVersion != "" {
			cmdToRunEnv = append(cmdToRunEnv, fmt.Sprintf("GOARM=%s", target.productArmVersion))
		}

		var debugFlags = ""
		if false {
			debugFlags = "-tags debug"
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
		if err := cmdToRun.Run(); err != nil {
			exitIfErr(fmt.Errorf("build %q: %w", platform, err))
		}
	}

	exitIfErr(compressOutputs(productOutputFolder, buildConfig.Config))
}

type buildTarget struct {
	productOS      string
	productOSLabel string

	productArch      string
	productArchLabel string

	productArmVersion string

	productArchWithDetails string

	productExt string
}

func parsePlatform(platform string) (buildTarget, error) {
	platform = strings.TrimSpace(platform)

	target := buildTarget{}

	parts := strings.Split(platform, "/")
	if len(parts) < 2 || len(parts) > 3 {
		return target, fmt.Errorf("invalid platform %q: expected GOOS/GOARCH or GOOS/GOARCH/GOARM", platform)
	}

	// OS
	target.productOS = parts[0]

	target.productOSLabel = target.productOS
	if display, ok := osToDisplay[target.productOS]; ok {
		target.productOSLabel = display
	}

	// ARCH
	target.productArch = parts[1]
	if target.productOS == "" || target.productArch == "" {
		return target, fmt.Errorf("invalid platform %q: GOOS and GOARCH are required", platform)
	}

	target.productArchLabel = target.productArch
	if display, ok := archToDisplay[target.productArch]; ok {
		target.productArchLabel = display
	}

	// ARM
	if len(parts) == 3 {
		if parts[2] == "" {
			return target, fmt.Errorf("invalid platform %q: GOARM is required when a third segment is provided", platform)
		}
		if target.productArch != "arm" {
			return target, fmt.Errorf("invalid platform %q: GOARM can only be set for arm builds", platform)
		}

		target.productArmVersion = parts[2]
	}

	// DET
	target.productArchWithDetails = target.productArchLabel
	if target.productArmVersion != "" {
		target.productArchWithDetails = fmt.Sprintf("%s-v%s", target.productArchLabel, target.productArmVersion)
	}

	// EXT
	target.productExt = osToExt[target.productOS]

	return target, nil
}

func compressOutputs(productOutputFolder string, config Config) error {
	if !config.Compress {
		return nil
	}

	compressorPath, err := exec.LookPath(config.Compressor)
	if err != nil {
		log.Printf("compress : skipping, %s not found in PATH\n", config.Compressor)
		return nil
	}

	return filepath.WalkDir(productOutputFolder, func(path string, dirEntry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if dirEntry.IsDir() {
			return nil
		}

		preSize := getFileSize(path)
		preSizeKB := preSize / 1024
		preSizeMB := preSize / (1024 * 1024)

		var stderrBuffer bytes.Buffer

		cmdToRun := exec.Command(compressorPath)
		if len(config.CompressorFlags) > 0 {
			cmdToRun.Args = append(cmdToRun.Args, config.CompressorFlags...)
		}
		cmdToRun.Args = append(cmdToRun.Args, path)
		cmdToRun.Stderr = io.MultiWriter(&stderrBuffer) //os.Stderr
		// cmdToRun.Stdout = os.Stdout
		if err := cmdToRun.Run(); err != nil {
			subErr := stderrBuffer.String()
			subErrI := strings.Index(strings.ToLower(subErr), "exception") + len("exception")
			if subErrI < len(subErr) {
				subErr = subErr[subErrI:]
			}
			return fmt.Errorf("compress : SKIP %s -> %s -> %s", dirEntry.Name(), config.Compressor, subErr)
		}

		postSize := getFileSize(path)
		postSizeKB := postSize / 1024
		postSizeMB := postSize / (1024 * 1024)

		log.Printf("compress : %s %d/%dK/%dMB -> %d/%dK/%dMB \n", dirEntry.Name(), preSize, preSizeKB, preSizeMB, postSize, postSizeKB, postSizeMB)

		return nil
	})
}
