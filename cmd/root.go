package cmd

import "github.com/spf13/cobra"

const version = "1.0.2"

var rootCmd = &cobra.Command{
	Use:     "crud",
	Version: version,
	Short:   "Crud code generation tool",
	Long:    "This is a tool that quickly generates CRUD interfaces based on database table structures.",
	Example: "crud init|gen",
}

func init() {
	// 根据openapi 生成接口文档
	rootCmd.AddCommand(initCmd, reloadCmd, startCmd, upgradeCmd)
}

func Run() error {
	return rootCmd.Execute()
}
