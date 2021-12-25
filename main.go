package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	"wfRadius/model"
	"wfRadius/src/request"
	"wfRadius/src/wifilog"
	"wfRadius/storage"
	"wfRadius/util"
	"wfRadius/ws"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"layeh.com/radius"
)

/*
TODO 将用户登录记录发回线上系统
*/
func prog(state overseer.State) {
	var err error

	// 配置日志模块
	var logger *zap.Logger
	if util.Config.Environment == model.EnvirIsProd {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename: os.Args[1] + "/zaplog",
			MaxSize:  10, // megabytes
		})
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			w,
			zap.InfoLevel,
		)
		logger = zap.New(core)
	} else if util.Config.Environment == model.EnvirIsDev {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(fmt.Sprintf("创建日志模块出错: %s\n", err.Error()))
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// 输出当前版本和构建时间
	zap.S().Infof("GitTag: %s", util.GitTag)
	zap.S().Infof("BuildTime: %s", util.BuildTime)

	// 配置数据库
	storage.Init()
	// 配置Http请求
	request.Init(util.Config.Token, util.Config.HTTPBackend)

	go wifilog.BeginUploadTask()
	go ws.Start(util.Config)

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
	viper.AddConfigPath(os.Args[1])
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件出错: %s", err.Error()))
	}
	err = viper.Unmarshal(&util.Config)
	if err != nil {
		panic(fmt.Errorf("解析配置文件出错: %s", err.Error()))
	}

	var f fetcher.Interface
	if util.Config.Environment == model.EnvirIsDev {
		f = &fetcher.File{
			Path: "wfRadius.next",
		}
	} else if util.Config.Environment == model.EnvirIsProd {
		f = &fetcher.HTTP{
			URL:      fmt.Sprintf("http://file.atsuas.cn/wfRadius_%s_%s", runtime.GOOS, runtime.GOARCH),
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
