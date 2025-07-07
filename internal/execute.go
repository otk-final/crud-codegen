package internal

import (
	"bytes"
	"crud-codegen/schema"
	"crud-codegen/tmpl"
	"fmt"
	"github.com/samber/lo"
	"maps"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Env struct {
	Driver     string      `json:"driver,omitempty"`     //驱动
	Url        string      `json:"url,omitempty"`        //information_schema mysql地址
	Datasource string      `json:"datasource,omitempty"` //目标数据库
	Config     Config      `json:"config,omitempty"`     //默认配置
	Tables     []*Endpoint `json:"tables,omitempty"`     //目标数据库表
}

type Config struct {
	CamelCase bool                     `json:"camel_case,omitempty"`
	Types     map[string]string        `json:"types,omitempty"`   //类型重命名
	Alias     map[string]string        `json:"alias,omitempty"`   //属性重命名
	Outputs   map[string]schema.Output `json:"outputs,omitempty"` //输入路径
	Api       schema.Api               `json:"api,omitempty"`     //接口配置
	Inherit   schema.Inherit           `json:"inherit,omitempty"` //继承
	Rewrite   bool                     `json:"rewrite,omitempty"` //是否强制重新生成文件
}

type Endpoint struct {
	Module      string                   `json:"module,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Endpoint    string                   `json:"endpoint,omitempty"`
	Comment     string                   `json:"comment,omitempty"`
	TablePrefix string                   `json:"table_prefix,omitempty"` //表名前缀
	TableName   string                   `json:"table_name,omitempty"`   //表名
	Outputs     map[string]schema.Output `json:"outputs,omitempty"`      //输出
	Renders     []*schema.Render         `json:"renders,omitempty"`      //字段渲染
}

type Mode struct {
	*schema.Table
	Module        string
	Name          string
	Endpoint      string
	Comment       string
	Inherit       schema.Inherit        //继承类
	Api           schema.Api            //接口配置
	Output        schema.Output         //输出
	RenderColumns []schema.RenderColumn //字段渲染
}

func Execute(env *Env) error {

	// loop
	for _, endpoint := range env.Tables {

		if env.Driver == "mysql" {
			//TODO
		} else if env.Driver == "oracle" {
			//TODO
		} else if env.Driver == "postgresql" {
			//TODO
		}

		tb, err := mysqlQuery(env.Url, env.Datasource, endpoint.TableName)
		if err != nil {
			return err
		}

		//导出文件
		err = export(env.Config, tb, endpoint)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Generate SUCCESS\n")

	return nil
}

func export(setting Config, table *schema.Table, endpoint *Endpoint) error {

	//忽略前缀
	tableName := table.Name
	if endpoint.TablePrefix != "" {
		tableName = strings.TrimPrefix(tableName, endpoint.TablePrefix)
	}

	//别名
	alias := setting.Alias
	types := setting.Types

	// comment 解析
	for _, column := range table.Columns {

		//解析注释
		text, enums := schema.ColumnComment(column.Comment).Parse()
		column.Comment = text
		column.Enums = enums

		//属性名
		field, ok := alias[column.Name]
		if !ok {
			//不存在重命名 则判断是否驼峰转换
			field = lo.Ternary(setting.CamelCase, CamelCase(column.Name, false), column.Name)
		}
		column.Alias = lo.Ternary(ok, column.Name, field)

		//类型转换
		langType, ok := types[column.Type]
		column.TypeAlias = lo.Ternary(ok, langType, column.Type)

		//指定字段重命名
		fieldAlias, ok := types["~"+column.Alias]
		if ok {
			column.TypeAlias = fieldAlias
		}
		fieldName, ok := types["~"+column.Name]
		if ok {
			column.TypeAlias = fieldName
		}
	}

	//过滤 继承属性
	table.Columns = lo.Filter(table.Columns, func(item *schema.Column, index int) bool {
		return !lo.Contains(setting.Inherit.Columns, item.Name)
	})

	//merge
	defaultOutput := maps.Clone(setting.Outputs)
	endpointOutput := endpoint.Outputs
	for k, v := range endpointOutput {
		defaultOutput[k] = v
	}

	//驼峰
	tableName = lo.Ternary(setting.CamelCase, CamelCase(tableName, true), tableName)

	mode := &Mode{
		Table:    table,
		Module:   endpoint.Module,
		Endpoint: lo.Ternary(endpoint.Endpoint != "", endpoint.Endpoint, table.Name),
		Name:     lo.Ternary(endpoint.Name != "", endpoint.Name, tableName),
		Comment:  lo.Ternary(endpoint.Comment != "", endpoint.Comment, table.Comment),
		Inherit:  setting.Inherit,
		Api:      setting.Api,
	}

	renderColumns := make([]schema.RenderColumn, 0)
	if len(endpoint.Renders) > 0 {
		for _, item := range endpoint.Renders {
			//查询字段
			targetName := item.Name
			matchColumn, ok := lo.Find(table.Columns, func(column *schema.Column) bool {
				return column.Name == targetName || column.Alias == targetName
			})
			if !ok {
				continue
			}

			if item.Label == "" {
				item.Label = matchColumn.Comment
			}

			renderColumns = append(renderColumns, schema.RenderColumn{
				Column: matchColumn,
				Render: item,
			})
		}
		mode.RenderColumns = renderColumns
	} else {
		//默认text
		mode.RenderColumns = lo.Map(table.Columns, func(item *schema.Column, index int) schema.RenderColumn {
			return schema.RenderColumn{
				Column: item,
				Render: &schema.Render{
					Name:      item.Name,
					Label:     item.Comment,
					Component: "div",
				},
			}
		})
	}

	//替换占位符
	mode.Api.Path = mode.Format(mode.Api.Path)

	//生成文件
	for name, output := range defaultOutput {

		if output.Ignore || output.File == "none" || output.File == "" {
			continue
		}

		//格式化文件路径
		file := mode.Format(output.File)
		//文件存在默认不覆盖
		_, err := os.Stat(file)
		//用户强制覆盖
		rewrite := os.IsNotExist(err) || output.Rewrite || setting.Rewrite
		if !rewrite {
			continue
		}

		//根据模版生成文件
		tp, err := tmpl.New(name, output.Template, output.Variables)
		if err != nil {
			return err
		}

		//格式化
		output.Header = lo.Map(output.Header, func(item string, index int) string {
			return mode.Format(item)
		})
		mode.Output = output

		//调用模版
		buf := &bytes.Buffer{}
		err = tp.Execute(buf, mode)
		if err != nil {
			return err
		}

		//创建目录
		file = tmpl.PwdJoinPath(file)
		dir := filepath.Dir(file)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		//写入文件
		_ = os.WriteFile(file, buf.Bytes(), os.ModePerm)

		fmt.Printf("Output：%s \n", file)
	}

	return nil
}

func (receiver *Mode) Format(text string) string {
	text = strings.ReplaceAll(text, "{module}", receiver.Module)
	text = strings.ReplaceAll(text, "{name}", receiver.Name)
	text = strings.ReplaceAll(text, "{table_name}", receiver.Table.Name)
	return text
}

// CamelCase 转换为 CamelCase 或 camelCase
func CamelCase(s string, upperFirst bool) string {
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.Title(parts[i])
	}
	result := strings.Join(parts, "")
	if !upperFirst && len(result) > 0 {
		runes := []rune(result)
		runes[0] = unicode.ToLower(runes[0])
		result = string(runes)
	}
	return result
}
