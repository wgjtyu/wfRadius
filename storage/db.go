package storage

import (
	"fmt"
	"os"
	"wfRadius/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB 数据库
var DB *gorm.DB

// Begin 开启Transaction
func Begin() *gorm.DB {
	return DB.Begin()
}

// Init 初始化数据库
func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open(os.Args[1] + "/main.db"))
	if err != nil {
		panic(fmt.Errorf("初始化数据库出错: %s", err.Error()))
	}
	_ = DB.AutoMigrate(&model.MWifiCode{})
	_ = DB.AutoMigrate(&model.MWifiLog{})
}
