package main

import (
	"fmt"
	"net/http"

	"github.com/imroc/req"
	"gorm.io/gorm"
)

// LoadData 从后端加载用户wifi数据
func LoadData() {
	r := req.New()
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = config.Token
	resp, err := r.Get(config.Backend+"/api/wifi_code/list", cookie)
	if err != nil {
		fmt.Printf("出错了%s\n", err.Error())
		return
	}

	var codes []MWifiCode
	var result *gorm.DB
	resp.ToJSON(&codes)
	if len(codes) > 0 {
		tx := db.Begin()
		for _, c := range codes {
			fmt.Printf("保存数据: %v\n", c)
			result = tx.Create(&c)
			if result.Error != nil {
				fmt.Printf("保存WifiCodes数据出错: %s\n", result.Error.Error())
				tx.Rollback()
				return
			}
		}
		tx.Commit()
	}

	fmt.Printf("%v\n", codes)
}
