package cmd

import "github.com/spf13/cobra"

const version = "1.0.4"

var rootCmd = &cobra.Command{
	Use:     "crud",
	Version: version,
	Short:   "CRUD code generation tool",
	Long:    "This is a tool that quickly generates CRUD interfaces based on database table structures.",
	Example: "crud init",
}

func init() {
	// 根据openapi 生成接口文档
	rootCmd.AddCommand(initCmd, reloadCmd, startCmd, upgradeCmd)
}

func Run() error {
	return rootCmd.Execute()
}
