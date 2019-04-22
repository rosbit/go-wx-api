package wxauth

import (
	"fmt"
	"time"
	"github.com/olebedev/config"
	"strings"
	"github.com/rosbit/go-wx-api/conf"
)

type WxUser struct {
	openId string
	accessToken string
	expireTime int64
	refreshToken string
	scope []string
	wxParams *wxconf.WxParamsT
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

var _userFields = []string{"nickname", "sex", "province", "city", "country", "headimgurl", "privilege", "unionid"}

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
	j, err := config.ParseJson(string(body))
	if err != nil {
		return err
	}
	if errcode, err := j.Int("errcode"); err != nil {
		errMsg := j.UString("errmsg", "")
		return fmt.Errorf("%v: %v", errcode, errMsg)
	}
	user.setUserInfo(j)
	return nil
}

func (user *WxUser) setUserInfo(res *config.Config) {
	for _, field := range _userFields {
		fmt.Printf("userInfo: %s", field)
	}
}

func (user *WxUser) getAccessToken(url string) error {
	res, err := CallWxAPI(url, "GET", nil)
	if err != nil {
		return err
	}
	fmt.Printf("get accessToken ok, res: %v\n", string(res))
	j, err := config.ParseJson(string(res))
	if err != nil {
		return err
	}
	if errcode, err := j.Int("errcode"); err == nil {
		errMsg := j.UString("errmsg", "")
		return fmt.Errorf("%v: %v", errcode, errMsg)
	}
	user.openId = j.UString("openid", "")
	user.accessToken = j.UString("access_token", "")
	user.expireTime = time.Now().Unix() + int64(j.UInt("expires_in")) - 10
	user.refreshToken = j.UString("refresh_token", "")
	user.scope = strings.Split(j.UString("scope", ""), ",")
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
