package cmd

import (
	"crud-codegen/internal"
	"crud-codegen/schema"
	"crud-codegen/tmpl"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Generate code based on the configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {

		//读取配置文件
		if envFile == "" {
			envFile = tmpl.PwdJoinPath(defaultEnvFileName)
		}
		data, err := os.ReadFile(envFile)
		if err != nil {
			fmt.Println(err)
			return err
		}

		var env internal.Env
		_ = json.Unmarshal(data, &env)

		//指定 已经配置的 table
		if filter != "" {
			match := lo.Filter(env.Tables, func(item *internal.Endpoint, index int) bool {
				return lo.Contains([]string{item.Name, item.Endpoint, item.TableName}, filter)
			})
			if len(match) == 0 {
				return fmt.Errorf("not found name")
			}
			env.Tables = match
		}

		//快捷生成
		if table != "" {
			env.Tables = []*internal.Endpoint{
				{
					TableName: table,
				},
			}
		}

		//指定 output
		if len(outputs) > 0 {

			matchOutputs := func(src map[string]schema.Output, matched []string) map[string]schema.Output {
				for key, _ := range src {
					if !lo.Contains(matched, key) {
						delete(src, key)
					}
				}
				return src
			}

			//删除其他数据
			env.Config.Outputs = matchOutputs(env.Config.Outputs, outputs)

			for _, table := range env.Tables {
				table.Outputs = matchOutputs(table.Outputs, outputs)
			}
		}

		//全局替换
		if rewrite {
			for _, output := range env.Config.Outputs {
				output.Rewrite = rewrite
			}
		}

		//执行
		return internal.Execute(&env)
	},
}

var envFile string
var filter string
var rewrite bool
var table string
var outputs []string
var defaultEnvFileName = "crud.json"

func init() {

	//reload
	reloadCmd.Flags().StringVarP(&envFile, "env", "e", "", "customize env file")
	reloadCmd.Flags().StringVarP(&filter, "filter", "f", "", "filter table_name or endpoint")
	reloadCmd.Flags().StringVarP(&table, "table", "t", "", "customize table")
	reloadCmd.Flags().StringSliceVarP(&outputs, "output", "o", []string{}, "filter outputs")
	reloadCmd.Flags().BoolVarP(&rewrite, "rewrite", "r", false, "rewrite outputs file")
}
