package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/auth"
)

func GetUserInfo(name, openId string) (map[string]interface{}, error) {
	return wxauth.GetUserInfo(name, openId)
}
