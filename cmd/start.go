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

		if len(tables) == 0 {
			return errors.New("required table names")
		}

		defaultEnv.Tables = lo.Map(tables, func(item string, index int) *internal.TableEndpoint {
			return &internal.TableEndpoint{
				Module:    "demo",
				TableName: item,
			}
		})

		//执行
		return internal.Execute(defaultEnv)
	},
	Example: "crud start -s username:password@tcp(localhost:3306)/information_schema -d demo -t table_name",
}

var tables []string

func init() {

	//reload
	startCmd.Flags().StringVarP(&defaultEnv.Url, "scheme datasource", "s", "", "scheme datasource address")
	startCmd.Flags().StringVarP(&defaultEnv.Datasource, "datasource", "d", "", "table datasource name")
	startCmd.Flags().StringSliceVarP(&tables, "tables", "t", []string{}, "table name")
}
