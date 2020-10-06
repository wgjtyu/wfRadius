package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v2"
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

	fmt.Printf("Code: %v, ip: %v, Port: %v, srv: %v, srvType: %v\n", r.Code, ip, port, srv, srvType)
	fmt.Printf("UserName: %s, Password: %s\n", username, password)
	fmt.Printf("LocalAddr: %v\n", r.LocalAddr.String())
	fmt.Printf("RemoteAddr: %v\n", r.RemoteAddr.String())

	var wifiCode MWifiCode
	err := db.View(func(txn *badger.Txn) error {
		var key bytes.Buffer
		key.WriteString("ID")
		key.WriteString(username)
		fmt.Printf("Key=%s\n", key.String())
		item, err := txn.Get(key.Bytes())
		if err != nil {
			fmt.Printf("handler-txn.Get出错: %s\n", err.Error())
			return err
		}

		var valCopy []byte
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			fmt.Printf("handler-item.ValueCopy出错: %s\n", err.Error())
			return err
		}
		err = json.Unmarshal(valCopy, &wifiCode)
		if err != nil {
			fmt.Printf("handler-json.Unmarshal出错: %s\n", err.Error())
			return err
		}

		return nil
	})
	var code radius.Code
	if err == badger.ErrKeyNotFound {
		code = radius.CodeAccessReject
	} else if err != nil {
		fmt.Printf("db获取数据出错: %s\n", err.Error())
		return
	} else {
		fmt.Printf("wifiCode=%v\n", wifiCode)
		if wifiCode.Valid && wifiCode.WifiCode == password {
			fmt.Printf("wifiCode==password\n")
			code = radius.CodeAccessAccept
			// } else if srvType == 0 {
			// fmt.Printf("srvType==0\n")
			// code = radius.CodeAccessAccept
		} else {
			fmt.Printf("reject\n")
			code = radius.CodeAccessReject
		}
	}

	log.Printf("Writing %v to %v", code, r.RemoteAddr)
	w.Write(r.Response(code))
}
