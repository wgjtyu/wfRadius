package ws

import (
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"wfRadius/model"
	"wfRadius/util"
)


// msgPack 服务器发来消息的包结构
type msgPack struct {
	Tag    string      `json:"tag"`
	RawMsg interface{} `json:"rawMsg"`
}

func reader(c *websocket.Conn, ch chan bool) {
	zap.S().Debugf("reader: 启动")
	defer func() {
		zap.S().Debugf("reader: 关闭")
		ch <- true
	}()
	for {
		var result msgPack
		err := c.ReadJSON(&result)
		if err != nil {
			zap.L().Debug("reader: 读取JSON出错", zap.String("error", err.Error()))
			break
		}
		if result.Tag == "WIFI_CODE_ADD" || result.Tag == "WIFI_CODE_UPDATE" { // 添加或更新
			var code model.MWifiCode
			mapstructure.Decode(result.RawMsg, &code)
			_ = util.UpdateCode(&code)
		} else if result.Tag == "WIFI_CODE_GET_ALL" { // 获取所有
			var codes []model.MWifiCode
			_ = mapstructure.Decode(result.RawMsg, &codes)
			_ = util.SaveCodes(codes)
		} else if result.Tag == "COMMAND" { // 执行指令
			zap.L().Info("reader-获取到指令")
			var command util.CommandMessage
			mapstructure.Decode(result.RawMsg, &command)
			go util.Proceed(command)
		}
	}
}

