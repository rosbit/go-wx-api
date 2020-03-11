package wxauth

import (
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"github.com/rosbit/go-wx-api/conf"
)

type WxUserInfo map[string]interface{} // 由于文档上sex是string类型，实际是整型。干脆不用struct解析了

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

func GetUserInfo(accessToken, openId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", accessToken, openId)
	resp, err := CallWxAPI(url, "GET", nil)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	if err = json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	if errcode, ok := res["errcode"]; ok {
		errmsg, _ := res["errmsg"]
		return nil, fmt.Errorf("%v: %v", errcode, errmsg)
	}
	return res, nil
}

// get user info by common access token
// please call this method after calling getOpenId().
// this calling will succeed if params are valid.
func (user *WxUser) GetInfoByAccessToken() (map[string]interface{}, error) {
	token := NewAccessTokenWithParams(user.wxParams)
	accessToken, err := token.Get()
	if err != nil {
		return nil, err
	}
	return GetUserInfo(accessToken, user.openId)
}

// get user info with OAuth2 API.
// please call this method after calling getOpenId().
// this calling maybe fail if unauthorized
func (user *WxUser) GetInfo() error {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", user.accessToken, user.openId)
	body, err := CallWxAPI(url, "GET", nil)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &user.UserInfo); err != nil {
		return err
	}
	if errcode, ok := user.UserInfo["errcode"]; ok {
		errmsg, _ := user.UserInfo["errmsg"]
		return fmt.Errorf("%d: %s", errcode, errmsg)
	}

	return nil
}

func (user *WxUser) getAccessToken(url string) error {
	res, err := CallWxAPI(url, "GET", nil)
	if err != nil {
		return err
	}
	// fmt.Printf("get accessToken ok, res: %v\n", string(res))

	var token struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenId       string `json:"openid"`
		Scope        string `json:"scope"`
		Errcode      int    `json:"errcode,omitempty"`
		Errmsg       string `json:"errmsg,omitempty"`
	}
	if err = json.Unmarshal(res, &token); err != nil {
		return err
	}
	if token.Errcode > 0 {
		return fmt.Errorf("%d: %s", token.Errcode, token.Errmsg)
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
