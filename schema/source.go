package schema

import (
	"fmt"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

type Table struct {
	Name    string    `json:"name"`    //表名
	Comment string    `json:"comment"` //表名称
	Columns []*Column `json:"columns"` //字段集
}

type Column struct {
	PK        bool              `json:"pk"`
	Name      string            `json:"name"`       //表字段
	Alias     string            `json:"alias"`      //表字段-别名
	Type      string            `json:"type"`       //字段类型
	TypeAlias string            `json:"type_alias"` //字段类型-别名
	Attribute map[string]string `json:"attribute"`  //扩展数据
	Comment   string            `json:"comment"`    //备注
	Required  bool              `json:"required"`   //必填
	Enums     []Enum            `json:"enums"`      //枚举值
}

type Enum struct {
	Name  string `json:"name"`  //名称
	Value string `json:"value"` //值
}

type ColumnComment string

var commentReg = regexp.MustCompile(`\[(.*?)\]`)

func (receiver ColumnComment) Parse() (string, []Enum) {
	enum := make([]Enum, 0)
	text := string(receiver)

	//example 解锁方式 [ad:广告,vip:会员,pay:支付]
	matches := commentReg.FindStringSubmatch(text)
	if len(matches) > 1 {
		enumTexts := strings.Split(matches[1], ",")
		enum = lo.Map(enumTexts, func(item string, index int) Enum {
			arr := strings.Split(item, ":")
			if len(arr) == 1 {
				arr = append(arr, item)
			}
			return Enum{
				Name:  strings.TrimSpace(arr[1]),
				Value: strings.TrimSpace(arr[0]),
			}
		})
		text = text[strings.Index(text, "["):]
		return strings.TrimSpace(text), enum
	}
	return strings.TrimSpace(text), enum
}

func (receiver Column) Format() string {
	enums := lo.Map(receiver.Enums, func(item Enum, index int) string {
		return fmt.Sprintf("%s:%s", item.Value, item.Name)
	})
	if len(enums) == 0 {
		return receiver.Comment
	}
	return fmt.Sprintf("%s[%s]", receiver.Comment, strings.Join(enums, ","))
}
