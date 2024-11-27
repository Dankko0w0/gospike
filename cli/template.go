package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func initTemplateCmd() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Manage project templates",
		Long:  `List, add, or remove project templates`,
	}

	// 添加子命令
	templateCmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List available templates",
			Run: func(cmd *cobra.Command, args []string) {
				listTemplates()
			},
		},
		&cobra.Command{
			Use:   "add [template-name] [template-path]",
			Short: "Add a new template",
			RunE:  addTemplate,
		},
	)

	rootCmd.AddCommand(templateCmd)
}

func listTemplates() {
	templates := []string{"basic", "web", "cli", "microservice"}
	fmt.Println("Available templates:")
	for _, t := range templates {
		fmt.Printf("- %s\n", t)
	}
}

func addTemplate(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("template name and path are required")
	}
	// TODO: 实现模板添加逻辑
	return nil
}
