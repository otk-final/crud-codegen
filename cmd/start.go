package cmd

import (
	"crud-codegen/internal"
	"errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Quickly start",
	RunE: func(cmd *cobra.Command, args []string) error {
		if schemeUrl == "" {
			return errors.New("required scheme datasource url")
		}
		if datasource == "" {
			return errors.New("required table datasource name")
		}
		if len(tables) == 0 {
			return errors.New("required table names")
		}

		defaultEnv.Tables = lo.Map(tables, func(item string, index int) *internal.Endpoint {
			return &internal.Endpoint{
				Module:    "demo",
				TableName: item,
			}
		})

		defaultEnv.Url = schemeUrl
		defaultEnv.Datasource = datasource

		//执行
		return internal.Execute(defaultEnv)
	},
	Example: "crud start -s username:password@tcp(localhost:3306)/information_schema -d demo -t table_name",
}

var schemeUrl string
var datasource string
var tables []string

func init() {

	//reload
	startCmd.Flags().StringVarP(&schemeUrl, "scheme datasource", "s", "", "scheme datasource address")
	startCmd.Flags().StringVarP(&datasource, "datasource", "d", "", "table datasource name")
	startCmd.Flags().StringSliceVarP(&tables, "tables", "t", []string{}, "table name")
}
