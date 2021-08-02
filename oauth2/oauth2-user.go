package wxoauth2

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/conf"
	"github.com/rosbit/go-wx-api/v2/auth"
	"fmt"
	"time"
	"strings"
)

type WxUser struct {
	openId string
	accessToken string
	expireTime int64
	refreshToken string
	scope []string
	wxParams *wxconf.WxParamT

	UserInfo *wxauth.WxUserInfo
}

func NewWxUser(name string) *WxUser {
	params := wxconf.GetWxParams(name)
	if params == nil {
		return nil
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

// get user info by common access token
// please call this method after calling getOpenId().
// this calling will succeed if params are valid.
func (user *WxUser) GetInfoByAccessToken() (*wxauth.WxUserInfo, error) {
	token := wxauth.NewAccessToken(user.wxParams)
	accessToken, err := token.Get()
	if err != nil {
		return nil, err
	}
	return wxauth.GetUserInfo(accessToken, user.openId)
}

// get user info with OAuth2 API.
// please call this method after calling getOpenId().
// this calling maybe fail if unauthorized
func (user *WxUser) GetInfo() error {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", user.accessToken, user.openId)
	var res struct {
		callwx.BaseResult
		wxauth.WxUserInfo
	}
	if _, err := callwx.CallWx(url, "GET", nil, nil, callwx.HttpCall, &res); err != nil {
		return err
	}
	fmt.Printf("oauth2 userInfo: %v\n", res.WxUserInfo)
	user.UserInfo = &res.WxUserInfo
	return nil
}

func (user *WxUser) getAccessToken(url string) error {
	var token struct {
		callwx.BaseResult
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenId       string `json:"openid"`
		Scope        string `json:"scope"`
	}
	if _, err := callwx.CallWx(url, "GET", nil, nil, callwx.HttpCall, &token); err != nil {
		return err
	}

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
