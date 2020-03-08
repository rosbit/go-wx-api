package wxauth

import (
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"github.com/rosbit/go-wx-api/conf"
)

type WxUserInfo struct {
	OpenId   string `json:"openid"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Province string `json:"province"`
	City     string `json:"city"`
	Country  string `json:"country"`
	HeadImgUrl string  `json:"headimgurl"`
	Privilege []string `json:"privilege"`
	UnionId   string   `json:"unionid"`
}

type WxUser struct {
	openId string
	accessToken string
	expireTime int64
	refreshToken string
	scope []string
	wxParams *wxconf.WxParamsT

	UserInfo WxUserInfo
}

func NewWxUser(params *wxconf.WxParamsT) *WxUser {
	if params == nil {
		return &WxUser{wxParams: &wxconf.WxParams}
	}
	return &WxUser{wxParams: params}
}

// 其实是authorize
func (user *WxUser) GetOpenId(code string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		user.wxParams.AppId, user.wxParams.AppSecret, code,
	)
	err := user.getAccessToken(url)
	if err != nil {
		return "", err
	}
	return user.openId, nil
}

/*
get user info with OAuth2 API.
please call this method after calling getOpenId().
*/
func (user *WxUser) GetInfo() error {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", user.accessToken, user.openId)
	body, err := CallWxAPI(url, "GET", nil)
	if err != nil {
		return err
	}
	var res map[string]interface{}
	if err = json.Unmarshal(body, &res); err != nil {
		return err
	}
	if errcode, ok := res["errcode"]; ok {
		errMsg, _ := res["errmsg"]
		return fmt.Errorf("%v: %v", errcode, errMsg)
	}

	json.Unmarshal(body, &user.UserInfo)
	return nil
}

func (user *WxUser) getAccessToken(url string) error {
	res, err := CallWxAPI(url, "GET", nil)
	if err != nil {
		return err
	}
	fmt.Printf("get accessToken ok, res: %v\n", string(res))

	var j map[string]interface{}
	if err = json.Unmarshal(res, &j); err != nil {
		return err
	}
	if errcode, ok := j["errcode"]; ok  {
		errMsg, _ := j["errmsg"]
		return fmt.Errorf("%v: %v", errcode, errMsg)
	}

	var token struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenId       string `json:"openid"`
		Scope        string `json:"scope"`
	}
	json.Unmarshal(res, &token)

	user.openId       = token.OpenId
	user.accessToken  = token.AccessToken
	user.expireTime   = time.Now().Unix() + int64(token.ExpiresIn) - 10
	user.refreshToken = token.RefreshToken
	user.scope        = strings.Split(token.Scope, ",")
	return nil
}

func (user *WxUser) TokenExpired() bool {
	return time.Now().Unix() > user.expireTime
}

func (user *WxUser) RefreshToken() error {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s",
		user.wxParams.AppId,
		user.refreshToken,
	)
	return user.getAccessToken(url)
}
