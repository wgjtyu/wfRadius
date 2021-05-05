package ws

import (
	"github.com/imroc/req"
	"go.uber.org/zap"
	"net/http"
	"wfRadius/model"
	"wfRadius/util"
)

// LoadData 从后端加载用户wifi数据
func LoadData(config model.MConfig) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = config.Token

	// 获取最新的WifiCode列表
	r := req.New()
	resp, err := r.Get(config.HTTPBackend+"/api/wifi_code/list", cookie)
	if err != nil {
		zap.L().Error("LoadData-获取WifiCode列表出错", zap.Error(err))
		return
	}
	var codes []model.MWifiCode
	resp.ToJSON(&codes)
	_ = util.SaveCodes(codes)
}
