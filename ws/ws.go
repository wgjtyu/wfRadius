package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
	"wfRadius/model"
)

const (
	writeWait = 30 * time.Second
)

func Start(config model.MConfig) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = config.Token

	header := http.Header{
		"Origin": []string{"https://www.atsuas.cn"},
		"Cookie": []string{cookie.String()},
	}

	url := config.WSBackend + "/api/ws"
	reconnectCh := make(chan bool, 0)
	defer close(reconnectCh)

	stopPingCh := make(chan bool, 0)
	defer close(stopPingCh)

	for {
		time.Sleep(5 * time.Second) // 睡眠5秒，以便c.Close对startPing产生作用
		c, _, err := websocket.DefaultDialer.Dial(url, header)
		if err != nil {
			zap.L().Error("ws.Start-连接ws服务出错", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		LoadData(config)
		start(c, reconnectCh, stopPingCh)

		select {
		case <-reconnectCh:
			zap.L().Info("ws.Start-将于5s后重新连接")
			stopPingCh <- true
			updatePeriod(false)
			c.Close()
		}
	}
}

func start(c *websocket.Conn, reconnectCh chan bool, stopPingCh chan bool) {
	go startPing(c, stopPingCh)
	go reader(c, reconnectCh)
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
