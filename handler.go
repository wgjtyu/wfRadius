package main

import (
	"fmt"
	"log"

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

	var count int64
	dbres := db.Table("m_wifi_code").Where("user_phone=? AND wifi_code=?", username, password).Count(&count)
	if dbres.Error != nil {
		fmt.Printf("匹配Wifi手机号和密码出错: %s\n", dbres.Error.Error())
		return
	}

	var code radius.Code

	fmt.Printf("UserName: %s, Password: %s\n", username, password)
	fmt.Printf("LocalAddr: %v\n", r.LocalAddr.String())
	fmt.Printf("RemoteAddr: %v\n", r.RemoteAddr.String())

	if count == 1 {
		code = radius.CodeAccessAccept
	} else if srvType == 0 {
		code = radius.CodeAccessAccept
	} else {
		code = radius.CodeAccessReject
	}
	log.Printf("Writing %v to %v", code, r.RemoteAddr)
	w.Write(r.Response(code))
}
