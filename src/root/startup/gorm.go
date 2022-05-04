package startup

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"os"
	"time"
	"wfRadius/model"
)

// InitGorm 初始化Gorm
func InitGorm(zl *zap.Logger) (*gorm.DB, error) {
	// 注意 表的主键需为ID而不是Id
	logger := zapgorm2.New(zl)

	logger.SlowThreshold = time.Second
	logger.IgnoreRecordNotFoundError = true

	db, err := gorm.Open(sqlite.Open(os.Args[1]+"/main.db"), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化数据库出错: %w", err)
	}
	_ = db.AutoMigrate(&model.MWifiCode{})
	_ = db.AutoMigrate(&model.MWifiLog{})
	return db, nil
}
