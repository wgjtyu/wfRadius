package model

import "time"

// MWifiCode 用户Wifi登录密码
type MWifiCode struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"userId"` // 用户ID
	WifiCode   string    `json:"wifiCode"`
	Valid      bool      `json:"valid"` // 通过此字段来判定该用户是否被禁止连接wifi
	UpdateTime time.Time `json:"updateTime"`
}
