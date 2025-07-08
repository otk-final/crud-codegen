package tmpl

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"text/template"
	"unicode"
)

//go:embed jdbc/api.tmpl
var jdbcApi string

//go:embed jdbc/entity.tmpl
var jdbcEntity string

//go:embed jdbc/persist.tmpl
var jdbcPersist string

//go:embed mybatis/api.tmpl
var mybatisApi string

//go:embed mybatis/entity.tmpl
var mybatisEntity string

//go:embed mybatis/persist.tmpl
var mybatisPersist string

//go:embed mybatis/service.tmpl
var mybatisService string

//go:embed mybatis/service_impl.tmpl
var mybatisServiceImpl string

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func New(name string, style string, variable map[string]string) (*template.Template, error) {

	tp := template.New("crud")

	var tmplFunc = template.FuncMap{
		"Capitalize": Capitalize,
		"Variable": func(key string) string {
			return variable[key]
		},
	}

	tp.Funcs(tmplFunc)

	if style != "" {
		style = PwdJoinPath(style)
		styleContent, err := os.ReadFile(style)
		if err != nil {
			return nil, err
		}
		return tp.Parse(string(styleContent))
	}

	switch name {
	case "jdbc-entity":
		return tp.Parse(jdbcEntity)
	case "jdbc-api":
		return tp.Parse(jdbcApi)
	case "jdbc-persist":
		return tp.Parse(jdbcPersist)
	case "mybatis-api":
		return tp.Parse(mybatisApi)
	case "mybatis-entity":
		return tp.Parse(mybatisEntity)
	case "mybatis-persist":
		return tp.Parse(mybatisPersist)
	case "mybatis-service":
		return tp.Parse(mybatisService)
	case "mybatis-service-impl":
		return tp.Parse(mybatisServiceImpl)
	}

	return nil, errors.New("invalid template name: " + name)
}

func PwdJoinPath(name string) string {
	if filepath.IsAbs(name) {
		return filepath.Clean(name)
	}
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, name)
}
