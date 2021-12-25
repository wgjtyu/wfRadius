package ws

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
	"wfRadius/config"
)

var lastSuccessPeriod time.Duration

// 更新心跳间隔
func updatePeriod(success bool) {
	if lastSuccessPeriod == 0 {
		lastSuccessPeriod = config.Instance.PingPeriod * 2
	}

	if success && config.Instance.PingPeriod == lastSuccessPeriod {
		return
	}

	lastPeriod := config.Instance.PingPeriod
	zap.L().Info("更新心跳间隔",
		zap.Bool("success", success),
		zap.Duration("pingPeriod", config.Instance.PingPeriod),
		zap.Duration("lastSuccessPeriod", lastSuccessPeriod))

	if success && config.Instance.PingPeriod != lastSuccessPeriod {
		lastSuccessPeriod = config.Instance.PingPeriod
		config.Instance.PingPeriod += config.Instance.PingPeriod / 10
	} else {
		if config.Instance.PingPeriod > lastSuccessPeriod {
			zap.L().Info("似乎找到最佳心跳间隔")
			config.Instance.PingPeriod = lastSuccessPeriod
		} else {
			config.Instance.PingPeriod /= 2
		}
	}

	if config.Instance.PingPeriod != lastPeriod {
		config.SetPingPeriod()
	}
	zap.L().Info("下次心跳间隔", zap.Duration("pingPeriod", config.Instance.PingPeriod))
}

func startPing(conn *websocket.Conn, stopPingCh <-chan bool) {
	zap.L().Info("startPing begin", zap.Duration("pingPeriod", config.Instance.PingPeriod))
	ticker := time.NewTicker(config.Instance.PingPeriod)

	for {
		select {
		case <-ticker.C: // 发送Ping
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
			if err != nil {
				ticker.Stop()
				return // 发送失败时退出程序
			} else {
				updatePeriod(true)
			}
			ticker.Reset(config.Instance.PingPeriod)
		case <-stopPingCh:
			zap.L().Info("执行stopPing")
			ticker.Stop()
			return
		}
	}
}
