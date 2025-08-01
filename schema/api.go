package schema

type Api struct {
	Class string `json:"class"` //通用结构类
	Path  string `json:"path"`  //路径前缀
}

type Inherit struct {
	Class   string   `json:"class"`   //继承类
	Columns []string `json:"columns"` //继承字段
}

type Render struct {
	Name      string `json:"name"`      //字段匹配名称
	Label     string `json:"label"`     //字段标签
	Component string `json:"component"` //组件模版
}

type Path struct {
	Path       string       `json:"path"`
	Comment    string       `json:"comment"`
	Parameters []*Parameter `json:"parameters"`
}

type Parameter struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Predicate string `json:"predicate"`
}

type Output struct {
	Ignore    bool              `json:"ignore,omitempty"`
	Header    []string          `json:"header,omitempty"`    //文件头信息
	File      string            `json:"file,omitempty"`      //文件地址
	Template  string            `json:"template,omitempty"`  //模版文件
	Variables map[string]string `json:"variables,omitempty"` //变量集
	Rewrite   bool              `json:"rewrite,omitempty"`   //重新生成文件
}
