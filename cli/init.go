package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	projectName string
	template    string
	moduleName  string
)

func initInitCmd() {
	initCmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Go project",
		Long:  `Initialize a new Go project with the specified structure and templates`,
		RunE:  runInit,
	}

	// 添加命令特定的标志
	initCmd.Flags().StringVarP(&template, "template", "t", "basic", "project template to use (basic, web, cli)")
	initCmd.Flags().StringVarP(&moduleName, "module", "m", "", "module name for go.mod")

	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("project name is required")
	}

	projectName = args[0]
	if moduleName == "" {
		moduleName = projectName
	}

	// 创建项目目录
	if err := createProjectStructure(); err != nil {
		return err
	}

	// 初始化 go.mod
	if err := initGoMod(); err != nil {
		return err
	}

	// 生成配置文件
	if err := generateConfig(); err != nil {
		return err
	}

	fmt.Printf("Successfully initialized project %s\n", projectName)
	return nil
}

func createProjectStructure() error {
	dirs := []string{
		"cmd",
		"internal",
		"internal/api",
		"internal/config",
		"internal/models",
		"internal/services",
		"internal/utils",
		"pkg",
		"scripts",
		"docs",
		"test",
	}

	for _, dir := range dirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	return nil
}

func initGoMod() error {
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectName
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateConfig() error {
	configContent := `app:
  name: ` + projectName + `
  version: 0.1.0
  env: development

server:
  host: localhost
  port: 8080

database:
  host: localhost
  port: 5432
  name: ` + projectName + `
  user: postgres
  password: postgres
`

	configPath := filepath.Join(projectName, "config.yaml")
	return os.WriteFile(configPath, []byte(configContent), 0644)
}
