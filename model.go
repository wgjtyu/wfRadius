package main

import (
	"time"
)

// MWifiCode 用户Wifi登录密码
type MWifiCode struct {
	ID         uint64    `gorm:"column:id" json:"id"`
	UserID     uint64    `gorm:"column:user_id" json:"userId"`       // 用户ID
	UserPhone  string    `gorm:"column:user_phone" json:"userPhone"` // 用户手机号
	WifiCode   string    `gorm:"column:wifi_code" json:"wifiCode"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

// TableName 返回MWifiCode的表名
func (MWifiCode) TableName() string { return "m_wifi_code" }

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
	Backend     string
	Token       string
	Environment ENVTYPE
}
