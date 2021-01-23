package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
	"wfRadius/model"
)

const (
	pingPeriod = 60 * time.Second
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
	ch := make(chan bool, 0)
	defer close(ch)
	for {
		c, _, err := websocket.DefaultDialer.Dial(url, header)
		if err != nil {
			zap.L().Error("ws.Start-连接ws服务出错", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		LoadData(config)
		start(c, ch)
		_, ok := <- ch
		if ok {
			c.Close()
		}
	}
}

func start(c *websocket.Conn, ch chan bool) {
	go startPing(c)
	go reader(c, ch)
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