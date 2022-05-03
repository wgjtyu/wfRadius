package ws

import (
	"wfRadius/model"
	"wfRadius/src/config"
	"wfRadius/src/request"
	"wfRadius/util"
)

// LoadData 从后端加载用户wifi数据
func LoadData(config config.MConfig) {
	resp := request.Get("/api/wifi_code/list")
	if resp != nil {
		var codes []model.MWifiCode
		resp.ToJSON(&codes)
		_ = util.SaveCodes(codes)
	}
}
