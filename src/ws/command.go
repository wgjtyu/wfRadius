package ws

import (
	"encoding/json"
	"net/http"
	"time"
	"wfRadius/src/config"
	"wfRadius/src/root/startup"
	"wfRadius/util"

	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type cmdProcessor struct {
	cfg    *config.MConfig
	logger *zap.Logger
}

func NewCmdProcessor(cfg *config.MConfig, l *zap.Logger) *cmdProcessor {
	return &cmdProcessor{
		cfg:    cfg,
		logger: l.Named("cmdProcessor"),
	}
}

// Proceed 处理设备指令并返回结果
func (c *cmdProcessor) Proceed(msg commandMessage, quitCh chan<- bool) {
	var inMsg innerMessage
	err := json.Unmarshal([]byte(msg.Content), &inMsg)
	if err != nil {
		c.logger.Error("解析cmd失败", zap.Error(err))
		return
	}
	c.logger.Info("处理设备指令", zap.String("command", inMsg.Name))
	if inMsg.Name == "VERSION" {
		c.putResult(msg.CommandID, &map[string]interface{}{
			"GitTag":    util.GitTag,
			"BuildTime": util.BuildTime})
	} else if inMsg.Name == "GET_CONFIG" {
		var mapObj map[string]interface{}
		err := mapstructure.Decode(startup.Instance, &mapObj)
		if err != nil {
			c.logger.Warn("Proceed-转换config出错", zap.Error(err))
			return
		}
		c.putResult(msg.CommandID, &mapObj)
	} else if inMsg.Name == "SET_CONFIG" {
		for key, value := range inMsg.Params {
			c.logger.Info("Proceed-设置配置", zap.String("key", key), zap.String("value", value))
			viper.Set(key, value)
		}
		err := viper.WriteConfig()
		if err != nil {
			c.logger.Warn("Proceed-保存viper配置失败", zap.Error(err))
			return
		}
		c.putResult(msg.CommandID, &map[string]interface{}{"status": "ok"})
	} else if inMsg.Name == "REBOOT" {
		c.putResult(msg.CommandID, &map[string]interface{}{"status": "ok"})
		quitCh <- true
	} else {
		c.putResult(msg.CommandID, &map[string]interface{}{"error": "unknown command"})
	}
}

// putResult 返回设备指令的执行结果给后端
func (c *cmdProcessor) putResult(commandID uint64, result *map[string]interface{}) {
	time.Sleep(5 * time.Second)
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = startup.Instance.Token

	content, err := json.Marshal(result)
	if err != nil {
		c.logger.Warn("设备指令结果转成string出错", zap.String("error", err.Error()))
		return
	}

	body := map[string]interface{}{
		"commandId": commandID,
		"content":   string(content),
	}

	r := req.New() // FIXME 用src/request下的方法
	resp, err := r.Post(startup.Instance.HTTPBackend+"/api/device/command_result_upload", cookie, req.BodyJSON(&body))
	if err != nil {
		c.logger.Error("上传指令结果, 本地出错", zap.String("error", err.Error()))
		return
	}
	if resp.Response().StatusCode != 200 {
		c.logger.Error("上传指令结果, 后端出错",
			zap.Int("statusCode", resp.Response().StatusCode),
			zap.String("body", resp.String()))
	}
}

// commandMessage 设备消息外围结构
type commandMessage struct {
	CommandID uint64 `json:"commandId"`
	Content   string `json:"content"`
}

// innerMessage 设备指令消息结构
type innerMessage struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}
