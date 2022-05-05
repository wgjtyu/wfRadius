package ws

import (
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
	"wfRadius/src/config"
)

const (
	writeWait = 10 * time.Second
)

var Set = wire.NewSet(
	NewWorker,
)

type Worker struct {
	reconnectCh chan bool
	cfg         *config.MConfig
	db          *gorm.DB
	logger      *zap.Logger
}

func NewWorker(c *config.MConfig, l *zap.Logger, db *gorm.DB) *Worker {
	return &Worker{
		reconnectCh: make(chan bool, 1), // 避免执行优雅退出时，reader结束时发生阻塞
		cfg:         c,
		db:          db,
		logger:      l.Named("ws.Worker"),
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = w.cfg.Token

	header := http.Header{
		"Origin": []string{"https://www.atsuas.cn"},
		"Cookie": []string{cookie.String()},
	}

	url := w.cfg.WSBackend + "/api/ws"
	reconnectCh := make(chan bool, 0)
	defer close(reconnectCh)

	stopPingCh := make(chan bool, 0)
	defer close(stopPingCh)

	for {
		time.Sleep(5 * time.Second) // 睡眠5秒，以便c.Close对startPing产生作用
		c, resp, err := websocket.DefaultDialer.Dial(url, header)
		zap.L().Info("websocket dial", zap.Int("status", resp.StatusCode))
		if err != nil {
			zap.L().Error("ws.Start-连接ws服务出错", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		w.LoadData()
		w.start(c, reconnectCh, stopPingCh)

		select {
		case <-reconnectCh:
			zap.L().Info("ws.Start-将于5s后重新连接")
			stopPingCh <- true
			updatePeriod(false)
			c.Close()
		}
	}
}

func (w *Worker) start(c *websocket.Conn, reconnectCh chan bool, stopPingCh chan bool) {
	go w.startPing(c, stopPingCh)
	go w.reader(c, reconnectCh)
	// 订阅WifiCode的添加和更新事件
	err := c.WriteJSON(map[string]interface{}{
		"type":   0, // 订阅
		"rawMsg": "WIFI_CODE_ADD",
	})
	if err != nil {
		zap.S().Errorf("LoadData-订阅添加WifiCode事件出错: %s", err.Error())
	}
	c.WriteJSON(map[string]interface{}{
		"type":   0, // 订阅
		"rawMsg": "WIFI_CODE_UPDATE",
	})
	if err != nil {
		zap.S().Errorf("LoadData-订阅更新WifiCode事件出错: %s", err.Error())
	}
}
