package main

import (
	"fmt"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/wgjtyu/goutil/overseer_fetcher"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"runtime"
	"time"
	"wfRadius/src/config"
	"wfRadius/src/request"
	"wfRadius/src/root/startup"
	"wfRadius/util"
)

func prog(state overseer.State) {
	rand.Seed(time.Now().UnixNano())
	app, err := BuildApp()
	if err != nil {
		zap.L().Error("启动出错", zap.Error(err))
		os.Exit(-1)
	}

	// 输出当前版本和构建时间
	zap.L().Info("Start",
		zap.String("GOOS", runtime.GOOS),
		zap.String("GOARCH", runtime.GOARCH),
		zap.String("GitTag", util.GitTag),
		zap.String("BuildTime", util.BuildTime))

	// 配置Http请求
	request.Init(startup.Instance.Token, startup.Instance.HTTPBackend)
	app.Run() // 在shutdown之前，会停在这里
	app.WaitForEnd()
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
