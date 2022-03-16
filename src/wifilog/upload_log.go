package wifilog

import (
	"time"
	"wfRadius/model"
	"wfRadius/src/request"
	"wfRadius/storage"

	"go.uber.org/zap"
)

func uploadLog() {
	var logs []model.MWifiLog
	err := storage.DB.Find(&logs).Error
	if len(logs) == 0 {
		return
	}
	if err != nil {
		zap.L().Error("uploadLog-获取所有wifi认证日志出错", zap.Error(err))
	} else {
		zap.L().Info("uploadLog-上传Wifi验证记录", zap.Int("count", len(logs)))
		resp := request.Post("/api/wifi_code/save_log", logs)

		// 删除已经上传的记录
		if resp != nil && resp.Response().StatusCode == 200 {
			var ids []uint64
			for _, l := range logs {
				ids = append(ids, l.ID)
			}
			err = storage.DB.Delete(&model.MWifiLog{}, ids).Error
			if err != nil {
				zap.L().Error("删除wifiLog出错", zap.Error(err))
			}
		}
	}
}

func BeginUploadTask() {
	d := time.Second * 300

	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		uploadLog()
	}
}
