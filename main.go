package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"layeh.com/radius"
)

var (
	// BuildTime 构建时间
	BuildTime = "Test Version"
	// GitTag Git的Tag标签
	GitTag = "Test Version"
)

var db *badger.DB
var config MConfig

/*
TODO 将用户登录记录发回线上系统
*/
func prog(state overseer.State) {
	var err error

	// 配置日志模块
	var logger *zap.Logger
	if config.Environment == EnvirIsProd {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename: "zaplog",
			MaxSize:  10, // megabytes
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.InfoLevel,
		)
		logger = zap.New(core)
	} else if config.Environment == EnvirIsDev {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(fmt.Sprintf("创建日志模块出错: %s\n", err.Error()))
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	zap.S().Infof("GitTag: %s", GitTag)
	zap.S().Infof("BuildTime: %s", BuildTime)

	db, err = badger.Open(badger.DefaultOptions("./db/"))
	if err != nil {
		zap.S().Errorf("连接数据库出错: %s", err.Error())
		return
	}

	c := cron.New()
	c.AddFunc("*/20 * * * *", func() {
		zap.L().Info(("重新加载数据"))
		LoadData()
	})
	c.Start()
	LoadData() // 首次重新加载数据

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	zap.L().Info("Starting server on :1812")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("退出程序")
}

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

	var f fetcher.Interface
	if config.Environment == EnvirIsDev {
		f = &fetcher.File{
			Path: "wfRadius.next",
		}
	} else if config.Environment == EnvirIsProd {
		f = &fetcher.HTTP{
			URL:      "http://file.atsuas.cn/wfRadius",
			Interval: 30 * time.Minute,
		}
	}
	overseer.Run(overseer.Config{
		Program:   prog,
		NoRestart: true,
		Address:   ":1812",
		Fetcher:   f,
	})
}
