package wifilog

import (
	"context"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"sync"
	"wfRadius/model"
	"wfRadius/src/request"
)

type Uploader struct {
	db      *gorm.DB
	limiter *rate.Limiter
	waitMux sync.Mutex
	logger  *zap.Logger
}

func NewUploader(db *gorm.DB, l *zap.Logger) *Uploader {
	u := &Uploader{
		db:      db,
		limiter: rate.NewLimiter(0.1, 1), // 1秒补充0.1个令牌
		logger:  l.Named("WifiLogUploader"),
	}
	return u
}

func (u *Uploader) uploadLog() {
	var logs []model.MWifiLog
	err := u.db.Find(&logs).Error
	if len(logs) == 0 {
		return
	}
	if err != nil {
		u.logger.Error("uploadLog-获取所有wifi认证日志出错", zap.Error(err))
	} else {
		u.logger.Info("uploadLog-上传Wifi验证记录", zap.Int("count", len(logs)))
		resp := request.Post("/api/wifi_code/save_log", logs)

		// 删除已经上传的记录
		if resp != nil && resp.Response().StatusCode == 200 {
			var ids []uint64
			for _, l := range logs {
				ids = append(ids, l.ID)
			}
			err = u.db.Delete(&model.MWifiLog{}, ids).Error
			if err != nil {
				u.logger.Error("删除wifiLog出错", zap.Error(err))
			}
		}
	}
}

func (u *Uploader) CheckAndUpload() {
	if !u.waitMux.TryLock() {
		return
	}
	err := u.limiter.Wait(context.Background())
	if err != nil {
		u.logger.Error("limiter.Wait failed", zap.Error(err))
	}
	u.waitMux.Unlock()
	u.uploadLog()
}
