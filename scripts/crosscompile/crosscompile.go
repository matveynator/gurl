// Command crosscompile builds the supported GURL release binaries.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type target struct {
	OperatingSystem string
	Architecture    string
}

var supportedTargets = []target{
	{OperatingSystem: "darwin", Architecture: "amd64"},
	{OperatingSystem: "darwin", Architecture: "arm64"},
	{OperatingSystem: "freebsd", Architecture: "amd64"},
	{OperatingSystem: "freebsd", Architecture: "arm64"},
	{OperatingSystem: "linux", Architecture: "amd64"},
	{OperatingSystem: "linux", Architecture: "arm64"},
	{OperatingSystem: "openbsd", Architecture: "amd64"},
	{OperatingSystem: "openbsd", Architecture: "arm64"},
	{OperatingSystem: "windows", Architecture: "amd64"},
	{OperatingSystem: "windows", Architecture: "arm64"},
}

func main() {
	outputDirectory := flag.String("output", "dist", "release output directory")
	releaseVersion := flag.String("version", "dev", "version embedded in binaries")
	flag.Parse()

	if err := os.MkdirAll(*outputDirectory, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "create output directory: %v\n", err)
		os.Exit(1)
	}
	for _, buildTarget := range supportedTargets {
		if err := build(*outputDirectory, *releaseVersion, buildTarget); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func build(outputDirectory string, releaseVersion string, buildTarget target) error {
	executableSuffix := ""
	if buildTarget.OperatingSystem == "windows" {
		executableSuffix = ".exe"
	}
	outputName := fmt.Sprintf("gurl_%s_%s%s", buildTarget.OperatingSystem, buildTarget.Architecture, executableSuffix)
	outputPath := filepath.Join(outputDirectory, outputName)
	linkerFlags := "-s -w -X main.version=" + releaseVersion

	// Target and output values only control a local release build initiated by the operator.
	command := exec.Command("go", "build", "-trimpath", "-ldflags", linkerFlags, "-o", outputPath, ".") // #nosec G204
	command.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS="+buildTarget.OperatingSystem,
		"GOARCH="+buildTarget.Architecture,
	)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("build %s/%s: %w: %s", buildTarget.OperatingSystem, buildTarget.Architecture, err, output)
	}
	fmt.Println(outputName)
	return nil
}
