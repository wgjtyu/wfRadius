package main

import (
	"time"
)

// MWifiCode 用户Wifi登录密码
type MWifiCode struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"userId"` // 用户ID
	WifiCode   string    `json:"wifiCode"`
	Valid      bool      `json:"valid"` // 通过此字段来判定该用户是否被禁止连接wifi
	UpdateTime time.Time `json:"updateTime"`
}

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
