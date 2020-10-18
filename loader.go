package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gorilla/websocket"
	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	pongWait = 60 * time.Second
)

// msgPack 服务器发来消息的包结构
type msgPack struct {
	Tag    string      `json:"tag"`
	RawMsg interface{} `json:"rawMsg"`
}

func reader(c *websocket.Conn) {
	zap.S().Debugf("reader: 启动")
	defer func() {
		zap.S().Debugf("reader: 关闭")
		// ws关闭后，退出程序，gRunner重新开启程序发起ws连接
		if c != nil {
			c.Close()
		}
		panic(nil)
	}()
	for {
		var result msgPack
		err := c.ReadJSON(&result)
		if err != nil { //FIXME 后端进程关闭导致的出错，需要重新连接?
			zap.L().Debug("reader: 读取JSON出错", zap.String("error", err.Error()))
			break
		}
		if result.Tag == "WIFI_CODE_ADD" || result.Tag == "WIFI_CODE_UPDATE" { // 添加或更新
			var code MWifiCode
			mapstructure.Decode(result.RawMsg, &code)
			updateCode(&code)
		} else if result.Tag == "WIFI_CODE_GET_ALL" { // 获取所有
			var codes []MWifiCode
			mapstructure.Decode(result.RawMsg, &codes)
			saveCodes(codes)
		} else if result.Tag == "COMMAND" { // 执行指令
			zap.L().Info("reader-获取到指令")
			var command CommandMessage
			mapstructure.Decode(result.RawMsg, &command)
			go Proceed(command)
		}
	}
}

func updateCode(code *MWifiCode) error {
	err := db.Update(func(txn *badger.Txn) error {
		jc, err := json.Marshal(code)
		if err != nil {
			zap.S().Debugf("addCode: WifiCode转换成json出错: %s", err.Error())
			return err
		}
		var key bytes.Buffer
		key.WriteString("ID")
		key.WriteString(strconv.FormatUint(code.UserID, 10))
		err = txn.Set(key.Bytes(), jc)
		if err != nil {
			zap.L().Warn("数据库保存WifiCodes数据出错", zap.String("error", err.Error()))
			return err
		}
		return nil
	})
	return err
}

func saveCodes(codes []MWifiCode) error {
	if len(codes) > 0 {
		err := db.Update(func(txn *badger.Txn) error {
			for _, c := range codes { // FIXME 保存失败时要rollback
				jc, err := json.Marshal(c)
				if err != nil {
					zap.L().Warn("WifiCodes数据转换成json出错", zap.String("error", err.Error()))
					return err
				}
				var key bytes.Buffer
				key.WriteString("ID")
				key.WriteString(strconv.FormatUint(c.UserID, 10))
				err = txn.Set(key.Bytes(), jc)
				if err != nil {
					zap.L().Warn("数据库保存WifiCodes数据出错", zap.String("error", err.Error()))
					return err
				}
			}
			return nil
		})
		return err
	}
	return nil
}

// LoadData 从后端加载用户wifi数据
func LoadData() {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = config.Token
	header := http.Header{
		"Origin": []string{"https://www.atsuas.cn"},
		"Cookie": []string{cookie.String()},
	}

	url := config.WSBackend + "/api/ws"
	c, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		zap.S().Errorf("LoadData-连接ws服务出错: %s", err.Error())
	}
	go reader(c)
	// 订阅WifiCode的添加和更新事件
	err = c.WriteJSON(map[string]interface{}{
		"type":   0, // 订阅
		"rawMsg": "WIFI_CODE_ADD",
	})
	if err != nil {
		zap.S().Errorf("LoadData-订阅添加WifiCode事件出错: %s", err.Error())
	}
	c.WriteJSON(map[string]interface{}{
		"type":   0, // 订阅
		"rawMsg": "WIFI_CODE_UPDATE",
	})
	if err != nil {
		zap.S().Errorf("LoadData-订阅更新WifiCode事件出错: %s", err.Error())
	}

	// 获取最新的WifiCode列表
	r := req.New()
	resp, err := r.Get(config.HTTPBackend+"/api/wifi_code/list", cookie)
	if err != nil {
		zap.S().Errorf("LoadData-获取WifiCode列表出错: %s\n", err.Error())
		return
	}
	var codes []MWifiCode
	resp.ToJSON(&codes)
	saveCodes(codes)
}
