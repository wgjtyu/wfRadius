package main

import (
	"encoding/json"
	"net/http"

	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// CommandMessage 设备消息外围结构
type CommandMessage struct {
	CommandID uint64 `json:"commandId"`
	Content   string `json:"content"`
}

// InnerMessage 设备指令消息结构
type InnerMessage struct {
	Name string `json:"name"`
}

// Proceed 处理设备指令并返回结果
func Proceed(msg CommandMessage) {
	var inMsg InnerMessage
	json.Unmarshal([]byte(msg.Content), &inMsg)
	zap.L().Info("处理设备指令", zap.String("command", inMsg.Name))
	if inMsg.Name == "VERSION" {
		putResult(msg.CommandID, &map[string]interface{}{
			"GitTag":    GitTag,
			"BuildTime": BuildTime})
	} else if inMsg.Name == "GET_CONFIG" {
		var mapObj map[string]interface{}
		err := mapstructure.Decode(config, &mapObj)
		if err != nil {
			zap.L().Warn("Proceed-转换config出错", zap.String("error", err.Error()))
			return
		}
		putResult(msg.CommandID, &mapObj)
	} else {
		putResult(msg.CommandID, &map[string]interface{}{"error": "unknown command"})
	}
}

// putResult 返回设备指令的执行结果给后端
func putResult(commandID uint64, result *map[string]interface{}) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = config.Token

	content, err := json.Marshal(result)
	if err != nil {
		zap.L().Warn("设备指令结果转成string出错", zap.String("error", err.Error()))
		return
	}

	body := map[string]interface{}{
		"commandId": commandID,
		"content":   string(content),
	}

	r := req.New()
	resp, err := r.Post(config.HTTPBackend+"/api/device/command_result_upload", cookie, req.BodyJSON(&body))
	if err != nil {
		zap.L().Error("上传指令结果, 本地出错", zap.String("error", err.Error()))
		return
	}
	if resp.Response().StatusCode != 200 {
		zap.L().Error("上传指令结果, 后端出错",
			zap.Int("statusCode", resp.Response().StatusCode),
			zap.String("body", resp.String()))
	}
}
