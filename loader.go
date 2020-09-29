package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/imroc/req"
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
	resp.ToJSON(&codes)
	if len(codes) > 0 {
		err := db.Update(func(txn *badger.Txn) error {
			for _, c := range codes { // FIXME 保存失败时要rollback
				fmt.Printf("保存数据: %v\n", c)
				jc, err := json.Marshal(c)
				if err != nil {
					fmt.Printf("WifiCodes数据转换成json出错: %s\n", err.Error())
					return err
				}
				var key bytes.Buffer
				key.WriteString("PHONE")
				key.WriteString(c.UserPhone)
				err = txn.Set(key.Bytes(), jc)
				if err != nil {
					fmt.Printf("数据库保存WifiCodes数据出错: %s\n", err.Error())
					return err
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}
