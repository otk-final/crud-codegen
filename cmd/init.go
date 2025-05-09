package cmd

import (
	"crud-codegen/internal"
	"crud-codegen/schema"
	"crud-codegen/tmpl"
	"encoding/json"
	"github.com/spf13/cobra"
	"os"
)

var api = schema.Output{
	Header: []string{
		"package com.demo.{module}.controller;",
		"",
		"import com.demo.ApiResult;",
		"import com.demo.{module}.entity.{name}Entity;",
		"import com.demo.{module}.repository.{name}Repository;",
	},
	File: "src/main/java/com/demo/{module}/controller/{name}Controller.java",
}

var entity = schema.Output{
	Header: []string{
		"package com.demo.{module}.entity;",
		"",
		"import com.demo.BaseEntity;",
	},
	File: "src/main/java/com/demo/{module}/entity/{name}Entity.java",
}
var persist = schema.Output{
	Header: []string{
		"package com.demo.{module}.repository;",
		"",
		"import com.demo.{module}.entity.{name}Entity;",
	},
	File: "src/main/java/com/demo/{module}/repository/{name}Repository.java",
}

var defaultConfig = internal.Config{
	CamelCase: true,
	Alias:     map[string]string{},
	Types: map[string]string{
		"bigint":   "Long",
		"int":      "Integer",
		"varchar":  "String",
		"json":     "JsonNode",
		"decimal":  "BigDecimal",
		"tinyint":  "Integer",
		"text":     "String",
		"datetime": "LocalDateTime",
		"date":     "LocalDate",
	},
	Outputs: map[string]schema.Output{

		"mybatis-api":     api,
		"mybatis-entity":  entity,
		"mybatis-persist": persist,

		//"jdbc-api":     api,
		//"jdbc-entity":  entity,
		//"jdbc-persist": persist,
	},
	Api: schema.Api{
		Class: "ApiResult",
		Path:  "/v1",
	},
	Inherit: schema.Inherit{
		Class: "BaseEntity",
		Columns: []string{
			"id",
			"created_at",
			"created_by",
			"updated_at",
			"updated_by",
			"del_flag",
		},
	},
}

var defaultEnv = &internal.Env{
	Driver:     "mysql",
	Url:        "username:password@tcp(localhost:3306)/information_schema",
	Datasource: "demo",
	Config:     defaultConfig,
	Tables: []*internal.TableEndpoint{
		{
			Module:    "test",
			TableName: "example",
			Comment:   "测试",
		},
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize environment variable configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		file := tmpl.PwdJoinPath(defaultEnvFileName)
		data, _ := json.MarshalIndent(defaultEnv, "", "   ")
		return os.WriteFile(file, data, os.ModePerm)
	},
}
