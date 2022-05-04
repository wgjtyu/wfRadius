package root

import (
	"fmt"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"wfRadius/src/config"
	"wfRadius/src/handler"
	"wfRadius/src/root/startup"
	"wfRadius/src/wifilog"
	"wfRadius/ws"
)

var Set = wire.NewSet(
	startup.Set,
	wire.Struct(new(App), "Config", "DB", "Logger"),
)

type App struct {
	Logger   *zap.Logger
	DB       *gorm.DB
	Config   *config.MConfig
	Worker   *ws.Worker
	RServer  *handler.RadiusServer
	Uploader *wifilog.Uploader
	wg       sync.WaitGroup
}

func (a *App) Run() {
	a.wg.Add(1)
	defer a.wg.Done()

	var innerWg sync.WaitGroup
	innerWg.Add(2)
	go a.Worker.Start(&innerWg)
	go a.RServer.Serve()
	innerWg.Wait()
}

func (a *App) Shutdown() {
	a.wg.Add(1)
	defer a.wg.Done()

	//a.Manager.Shutdown() // 关闭所有视频流
}

func (a *App) WaitForEnd() {
	a.wg.Wait()

	a.Logger.Debug("WaitForEnd end")
	if a.Config.Environment == config.EnvirIsProd {
		err := a.Logger.Sync()
		if err != nil {
			fmt.Println("zap sync failed --- ", err)
		}
	}
}
