package startup

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"time"
	"wfRadius/src/config"
)

// Instance 保存应用的配置
var Instance config.MConfig

func InitCfg() *config.MConfig {
	viper.SetConfigName("config")
	viper.AddConfigPath(os.Args[1])
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件出错: %s", err.Error()))
	}
	err = viper.Unmarshal(Instance)
	if err != nil {
		panic(fmt.Errorf("解析配置文件出错: %s", err.Error()))
	}

	if Instance.PingPeriod == 0 { // 默认心跳间隔为60秒
		Instance.PingPeriod = 60 * time.Second
	}
	return &Instance
}

func SetPingPeriod() {
	viper.Set("PingPeriod", Instance.PingPeriod)
	err := viper.WriteConfig()
	if err != nil {
		zap.L().Error("更新PingPeriod配置失败", zap.Error(err))
	}
}
