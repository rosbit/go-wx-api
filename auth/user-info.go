package wxauth

import (
	"github.com/rosbit/go-wget"
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"fmt"
)

type WxUserInfo struct {
	OpenId   string `json:"openid"`
	NickName string `json:"nickname"`
	Sex int8 `json:"sex"`
	Province string `json:"province"`
	City     string `json:"city"`
	Country  string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	Privilege []string `json:"privilege"`
	UnionId string `json:"unionid"`
}

func GetUserInfo(name, openId string) (*WxUserInfo, error) {
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
	return &res.WxUserInfo, nil
}
