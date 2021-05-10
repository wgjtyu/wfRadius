package request

import (
	"net/http"

	"github.com/imroc/req"
	"go.uber.org/zap"
)

var _cookie *http.Cookie
var _backend string

func Init(token string, backend string) {
	_cookie = new(http.Cookie)
	_cookie.Name = "token"
	_cookie.Value = token
	_backend = backend
}

func Get(url string) *req.Resp {
	// 获取最新的WifiCode列表
	resp, err := req.Get(_backend+url, _cookie)
	if err != nil {
		zap.L().Error("request.Req-Get请求出错", zap.Error(err))
		return nil
	}
	return resp
}

func Post(url string, data interface{}) *req.Resp {
	resp, err := req.Post(_backend+url, req.BodyJSON(data), _cookie)
	if err != nil {
		zap.L().Error("request.Req-Post请求出错", zap.Error(err))
		return nil
	}
	return resp
}
