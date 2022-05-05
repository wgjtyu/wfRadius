package wifilog

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
	"wfRadius/model"
	"wfRadius/src/request"
)

type Uploader struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUploader(db *gorm.DB, l *zap.Logger) *Uploader {
	u := &Uploader{db: db, logger: l.Named("WifiLogUploader")}
	go u.beginTask()
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

// FIXME 改造成RateLimiter, 并在Shutdown时全部上传
func (u *Uploader) beginTask() {
	d := time.Second * 300

	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		u.uploadLog()
	}
}

func (u *Uploader) Shutdown() {
	// FIXME implement it
}
