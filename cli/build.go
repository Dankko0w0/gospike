package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func initBuildCmd() {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build the project for distribution",
		Long:  `Build and compile the project for distribution`,
		RunE:  runBuild,
	}

	buildCmd.Flags().String("output", "dist", "output directory for built files")
	buildCmd.Flags().Bool("cross-compile", false, "enable cross-compilation")
	buildCmd.Flags().StringP("target", "t", "", "specific target to build (e.g., cmd/main.go)")

	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	crossCompile, _ := cmd.Flags().GetBool("cross-compile")
	target, _ := cmd.Flags().GetString("target")

	// 创建输出目录
	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if crossCompile {
		return crossCompileBuild(output, target)
	}
	return normalBuild(output, target)
}

func normalBuild(output string, target string) error {
	args := []string{"build"}
	if target != "" {
		args = append(args, target)
	}
	args = append(args, "-o", filepath.Join(output, "app"))

	buildCmd := exec.Command("go", args...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	return buildCmd.Run()
}

func crossCompileBuild(output string, target string) error {
	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"windows", "amd64"},
		{"darwin", "amd64"},
	}

	for _, platform := range platforms {
		env := os.Environ()
		env = append(env, fmt.Sprintf("GOOS=%s", platform.os))
		env = append(env, fmt.Sprintf("GOARCH=%s", platform.arch))

		outputFile := filepath.Join(output, fmt.Sprintf("app-%s-%s", platform.os, platform.arch))
		if platform.os == "windows" {
			outputFile += ".exe"
		}

		args := []string{"build"}
		if target != "" {
			args = append(args, target)
		}
		args = append(args, "-o", outputFile)

		buildCmd := exec.Command("go", args...)
		buildCmd.Env = env
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			return fmt.Errorf("failed to build for %s/%s: %w", platform.os, platform.arch, err)
		}
	}

	return nil
}
