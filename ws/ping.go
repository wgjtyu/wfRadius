package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

func startPing(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		zap.L().Debug("关闭Ping")
		ticker.Stop()
	}()
	for {
		select {
		case <- ticker.C: // 发送Ping
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait)); err != nil {
				zap.L().Warn("Ping出错", zap.Error(err))
				return // 发送失败时退出程序
			}
		}
	}
}