package main

import (
	"fmt"
	"github.com/wgjtyu/logMansion/lib"
	"log"
	"os"
	"runtime"
	"time"
	"wfRadius/config"
	"wfRadius/model"
	"wfRadius/src/request"
	"wfRadius/src/wifilog"
	"wfRadius/storage"
	"wfRadius/util"
	"wfRadius/ws"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

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
	if config.Instance.Environment == model.EnvirIsProd {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   os.Args[1] + "/zaplog",
			MaxSize:    2, // megabytes
			MaxBackups: 5,
		})
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			w,
			zap.InfoLevel,
		)
		logger = zap.New(core)
		lmCore := lib.NewCore(config.Instance.LogBackend, config.Instance.LogProjectID,
			config.Instance.LogKey, zapcore.InfoLevel)
		logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, lmCore)
		}))
	} else if config.Instance.Environment == model.EnvirIsDev {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(fmt.Sprintf("创建日志模块出错: %s\n", err.Error()))
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// 输出当前版本和构建时间
	zap.L().Info("Start",
		zap.String("GOOS", runtime.GOOS),
		zap.String("GOARCH", runtime.GOARCH),
		zap.String("GitTag", util.GitTag),
		zap.String("BuildTime", util.BuildTime))

	// 配置数据库
	storage.Init()
	// 配置Http请求
	request.Init(config.Instance.Token, config.Instance.HTTPBackend)

	go wifilog.BeginUploadTask()
	go ws.Start(config.Instance)

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("退出程序")
}

func main() {
	config.InitCfg()

	var f fetcher.Interface
	if config.Instance.Environment == model.EnvirIsDev {
		f = &fetcher.File{
			Path: "wfRadius.next",
		}
	} else if config.Instance.Environment == model.EnvirIsProd {
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
