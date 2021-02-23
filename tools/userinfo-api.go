package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/auth"
)

func GetUserInfo(name, openId string) (*wxauth.WxUserInfo, error) {
	return wxauth.GetUserInfo(name, openId)
}
