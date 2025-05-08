package internal

type ApiName string

const (
	ApiCreate ApiName = "create"
	ApiUpdate         = "update"
	ApiDelete         = "delete"
	ApiGet            = "get"
	ApiPage           = "page"
	ApiList           = "list"
	ApiSearch         = "search"
)

func NewApi(setting *Config, name ApiName) {

}
