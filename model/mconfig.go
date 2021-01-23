package model

// ENVTYPE 运行环境类型
type ENVTYPE uint32

const (
	// EnvirIsProd 生产环境
	EnvirIsProd ENVTYPE = 0
	// EnvirIsDev 开发环境
	EnvirIsDev ENVTYPE = 1
)

// MConfig 配置文件的结构
type MConfig struct {
	Name        string
	HTTPBackend string `mapstructure:"http_backend"`
	WSBackend   string `mapstructure:"ws_backend"`
	Token       string
	Environment ENVTYPE
}