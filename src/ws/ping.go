package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
	"wfRadius/src/root/startup"
)

var lastSuccessPeriod time.Duration

// 更新心跳间隔
func updatePeriod(success bool) {
	if lastSuccessPeriod == 0 {
		lastSuccessPeriod = startup.Instance.PingPeriod * 2
	}

	if success && startup.Instance.PingPeriod == lastSuccessPeriod {
		return
	}

	lastPeriod := startup.Instance.PingPeriod
	zap.L().Info("更新心跳间隔",
		zap.Bool("success", success),
		zap.Duration("pingPeriod", startup.Instance.PingPeriod),
		zap.Duration("lastSuccessPeriod", lastSuccessPeriod))

	if success && startup.Instance.PingPeriod != lastSuccessPeriod {
		lastSuccessPeriod = startup.Instance.PingPeriod
		startup.Instance.PingPeriod += startup.Instance.PingPeriod / 10
	} else {
		if startup.Instance.PingPeriod > lastSuccessPeriod {
			zap.L().Info("似乎找到最佳心跳间隔")
			startup.Instance.PingPeriod = lastSuccessPeriod
		} else {
			startup.Instance.PingPeriod /= 2
		}
	}

	if startup.Instance.PingPeriod != lastPeriod {
		startup.SetPingPeriod()
	}
	zap.L().Info("下次心跳间隔", zap.Duration("pingPeriod", startup.Instance.PingPeriod))
}

func (w *Worker) startPing(conn *websocket.Conn, ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()
	zap.L().Info("startPing begin", zap.Duration("pingPeriod", startup.Instance.PingPeriod))
	ticker := time.NewTicker(startup.Instance.PingPeriod)

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
			ticker.Reset(startup.Instance.PingPeriod)
		case <-ctx.Done():
			zap.L().Info("执行stopPing")
			ticker.Stop()
			return
		}
	}
}
