package main

import (
	"fmt"
	"github.com/wgjtyu/logMansion/lib"
	"log"
	"os"
	"runtime"
	"wfRadius/src/config"
	"wfRadius/src/request"
	"wfRadius/src/root/startup"
	"wfRadius/src/wifilog"
	"wfRadius/storage"
	"wfRadius/util"
	"wfRadius/ws"

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/wgjtyu/goutil/overseer_fetcher"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"layeh.com/radius"
)

func prog(state overseer.State) {
	var err error

	// 配置日志模块
	var logger *zap.Logger
	if startup.Instance.Environment == config.EnvirIsProd {
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
		lmCore := lib.NewCore(startup.Instance.LogBackend, startup.Instance.LogProjectID,
			startup.Instance.LogKey, zapcore.InfoLevel)
		logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, lmCore)
		}))
	} else if startup.Instance.Environment == config.EnvirIsDev {
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
	request.Init(startup.Instance.Token, startup.Instance.HTTPBackend)

	go wifilog.BeginUploadTask()
	go ws.Start(startup.Instance)

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
	cfg := startup.InitCfg()

	var f fetcher.Interface
	if cfg.Environment == config.EnvirIsDev {
		f = &fetcher.File{
			Path: "wfRadius.next",
		}
	} else if cfg.Environment == config.EnvirIsProd {
		f = &overseerfetcher.HttpFetcher{
			URL:           fmt.Sprintf("https://afile.atsuas.cn/file/wfRadius_%s_%s", runtime.GOOS, runtime.GOARCH),
			WorkDirectory: os.Args[1],
		}
	}
	overseer.Run(overseer.Config{
		Program:   prog,
		NoRestart: true,
		Address:   ":1812",
		Fetcher:   f,
	})
}
