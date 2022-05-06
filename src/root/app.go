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
	"wfRadius/src/ws"
)

var Set = wire.NewSet(
	startup.Set,
	wire.Struct(new(App), "Config", "DB", "Logger", "Worker", "RServer", "Uploader"),
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
	go a.RServer.Serve(&innerWg)
	innerWg.Wait()
	a.Logger.Debug("Run end")
}

func (a *App) Shutdown() {
	a.wg.Add(1)
	defer a.wg.Done()

	a.Worker.Shutdown()
	err := a.RServer.Shutdown()
	if err != nil {
		a.Logger.Error("RServer.Shutdown()出错", zap.Error(err))
	}
	a.Uploader.Shutdown()

	a.Logger.Debug("Shutdown结束")
}

func (a *App) WaitForEnd() {
	a.Logger.Debug("WaitForEnd begin")
	a.wg.Wait()

	a.Logger.Debug("WaitForEnd end")
	if a.Config.Environment == config.EnvirIsProd {
		err := a.Logger.Sync()
		if err != nil {
			fmt.Println("zap sync failed --- ", err)
		}
	}
}
