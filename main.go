package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"layeh.com/radius"
)

var db *gorm.DB
var config MConfig

/*
从线上获取用户的登录账号和密码
登录账号为用户的手机号
密码为线上系统生成的密码
将用户登录记录发回线上系统
*/
func main() {
	var err error

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件出错: %s", err.Error()))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("解析配置文件出错: %s", err.Error()))
	}

	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&MWifiCode{})

	LoadData()

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	log.Printf("Starting server on :1812")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
