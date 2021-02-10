package wxauth

import (
	"github.com/rosbit/go-wget"
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"fmt"
)

type WxUserInfo map[string]interface{} // 由于文档上sex是string类型，实际是整型。干脆不用struct解析了

func GetUserInfo(name, openId string) (map[string]interface{}, error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", accessToken, openId)
		return
	}

	var res struct {
		callwx.BaseResult
		WxUserInfo
	}
	if _, err := CallWx(name, genParams, "GET", wget.HttpCallJ, &res); err != nil {
		return nil, err
	}
	return res.WxUserInfo, nil
}
