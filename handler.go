package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v2"
	"go.uber.org/zap"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func handler(w radius.ResponseWriter, r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)
	ip := rfc2865.LoginIPHost_Get(r.Packet)
	port := rfc2865.LoginTCPPort_Get(r.Packet)
	srv := rfc2865.LoginService_Get(r.Packet)
	srvType := rfc2865.ServiceType_Get(r.Packet)

	zap.L().Info("handler-用户请求信息", zap.String("Code", r.Code.String()),
		zap.String("IP", ip.String()),
		zap.String("Port", port.String()),
		zap.String("Srv", srv.String()),
		zap.String("SrvType", srvType.String()),
		zap.String("UserId", username),
		zap.String("Password", password),
		zap.String("LocalAddr", r.LocalAddr.String()),
		zap.String("RemoteAddr", r.RemoteAddr.String()))

	var wifiCode MWifiCode
	err := db.View(func(txn *badger.Txn) error {
		var key bytes.Buffer
		key.WriteString("ID")
		key.WriteString(username)
		fmt.Printf("Key=%s\n", key.String())
		item, err := txn.Get(key.Bytes())
		if err != nil {
			zap.L().Warn("handler-txn.Get出错", zap.String("error", err.Error()))
			return err
		}

		var valCopy []byte
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			zap.L().Warn("handler-item.ValueCopy出错", zap.String("error", err.Error()))
			return err
		}
		err = json.Unmarshal(valCopy, &wifiCode)
		if err != nil {
			zap.L().Warn("handler-json.Unmarshal出错", zap.String("error", err.Error()))
			return err
		}

		return nil
	})
	var code radius.Code
	if err == badger.ErrKeyNotFound {
		code = radius.CodeAccessReject
	} else if err != nil {
		zap.L().Warn("handler-db获取数据出错", zap.String("error", err.Error()))
		return
	} else {
		if wifiCode.Valid && wifiCode.WifiCode == password {
			code = radius.CodeAccessAccept
			// } else if srvType == 0 {
			// fmt.Printf("srvType==0\n")
			// code = radius.CodeAccessAccept
		} else {
			code = radius.CodeAccessReject
		}
	}

	log.Printf("Writing %v to %v", code, r.RemoteAddr)
	w.Write(r.Response(code))
}
