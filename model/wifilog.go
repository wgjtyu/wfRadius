package model

import "time"

// MWifiLog 用户Wifi验证记录
type MWifiLog struct {
	ID      uint64    `json:"id"`
	UserID  uint64    `json:"userId"` // 用户ID
	MacAddr string    `json:"macAddr"`
	Time    time.Time `json:"time"`
}
