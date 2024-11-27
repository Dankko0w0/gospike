package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gospike",
	Short: "GoSpike - A modern Go project scaffolding tool",
	Long: `GoSpike is a powerful CLI tool for creating and managing Go projects.
It helps you quickly scaffold new projects with best practices and common patterns.`,
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 添加全局标志
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path")

	// 初始化所有子命令
	initInitCmd()
	initBuildCmd()
	initTemplateCmd()
}
