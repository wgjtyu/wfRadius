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
func (w *Worker) updatePeriod(success bool) {
	if lastSuccessPeriod == 0 {
		lastSuccessPeriod = w.cfg.PingPeriod * 2
	}

	if success && w.cfg.PingPeriod == lastSuccessPeriod {
		return
	}

	lastPeriod := w.cfg.PingPeriod
	w.logger.Info("更新心跳间隔",
		zap.Bool("success", success),
		zap.Duration("pingPeriod", w.cfg.PingPeriod),
		zap.Duration("lastSuccessPeriod", lastSuccessPeriod))

	if success && w.cfg.PingPeriod != lastSuccessPeriod {
		lastSuccessPeriod = w.cfg.PingPeriod
		w.cfg.PingPeriod += w.cfg.PingPeriod / 10
	} else {
		if w.cfg.PingPeriod > lastSuccessPeriod {
			w.logger.Info("似乎找到最佳心跳间隔")
			w.cfg.PingPeriod = lastSuccessPeriod
		} else {
			w.cfg.PingPeriod /= 2
		}
	}

	if w.cfg.PingPeriod != lastPeriod {
		startup.SetPingPeriod()
	}
	w.logger.Info("下次心跳间隔", zap.Duration("pingPeriod", w.cfg.PingPeriod))
}

func (w *Worker) startPing(conn *websocket.Conn, ctx context.Context) {
	w.wg.Add(1)
	defer w.wg.Done()
	w.logger.Info("startPing begin", zap.Duration("pingPeriod", w.cfg.PingPeriod))
	ticker := time.NewTicker(w.cfg.PingPeriod)

	for {
		select {
		case <-ticker.C: // 发送Ping
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
			if err != nil {
				ticker.Stop()
				return // 发送失败时退出程序
			} else {
				w.updatePeriod(true)
				ticker.Reset(w.cfg.PingPeriod)
			}
		case <-ctx.Done():
			w.logger.Info("执行stopPing")
			ticker.Stop()
			return
		}
	}
}
