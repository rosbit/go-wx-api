package wxtools

import (
	"fmt"
	"encoding/json"
	"github.com/rosbit/go-wx-api/auth"
)

func GetUserInfo(accessToken, openId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", accessToken, openId)
	resp, err := wxauth.CallWxAPI(url, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	if err = json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	if _, ok := res["errcode"]; ok {
		return nil, fmt.Errorf("%s", string(resp))
	}
	return res, nil
}
