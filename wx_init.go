package wxapi

import (
	"github.com/rosbit/go-wx-api/v2/conf"
)

func InitWx(tokenStorePath string) {
	wxconf.InitTokenStorePath(tokenStorePath)
}

func SetWxParams(serviceName string, token, appId, appSecret, aesKey string) error {
	return wxconf.NewWxParams(serviceName, token, appId, appSecret, aesKey)
}
