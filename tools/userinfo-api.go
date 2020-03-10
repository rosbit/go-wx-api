package wxtools

import (
	"github.com/rosbit/go-wx-api/auth"
)

func GetUserInfo(accessToken, openId string) (map[string]interface{}, error) {
	return wxauth.GetUserInfo(accessToken, openId)
}
