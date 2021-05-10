package util

import (
	"wfRadius/model"
	"wfRadius/storage"

	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

func SaveCodes(codes []model.MWifiCode) error {
	if len(codes) == 0 {
		return nil
	}
	tx := storage.Begin()
	for _, code := range codes {
		res := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&code)
		if res.Error != nil {
			zap.L().Warn("数据库创建WifiCode数据出错", zap.Error(res.Error))
		}
	}
	tx.Commit()
	return nil
}
func UpdateCode(code *model.MWifiCode) error {
	tx := storage.Begin()
	res := tx.Save(code)
	if res.Error != nil {
		zap.L().Warn("数据库保存WifiCode数据出错", zap.Error(res.Error))
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}
