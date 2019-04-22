package wxauth

import (
	"time"
	"github.com/olebedev/config"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/rosbit/go-wx-api/conf"
)

type AccessToken struct {
	accessToken string
	expireTime int64
	wxParams *wxconf.WxParamsT
}

func NewAccessToken() *AccessToken {
	return NewAccessTokenWithParams(nil)
}

func NewAccessTokenWithParams(params *wxconf.WxParamsT) *AccessToken {
	if params == nil {
		params = &wxconf.WxParams
	}
	token := &AccessToken{wxParams:params}
	token.loadFromStore()
	return token
}

func (token *AccessToken) Get() (string, error) {
	if token.expired() {
		err := token.get_access_token()
		if err != nil {
			return "", err
		}
	}
	return token.accessToken, nil
}

func (token *AccessToken) expired() bool {
	return token.expireTime < time.Now().Unix()
}

func (token *AccessToken) get_access_token() error {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		token.wxParams.AppId,
		token.wxParams.AppSecret,
	)
	body, err := CallWxAPI(url, "GET", nil)
	if err != nil {
		return err
	}
	res, err := config.ParseJson(string(body))
	if err != nil {
		return err
	}
	token.accessToken = res.UString("access_token", "")
	token.expireTime = int64(res.UInt("expires_in", 0)) + time.Now().Unix() - 10

	return token.saveToStore()
}

func (token *AccessToken) savePath() string {
	return fmt.Sprintf("%s/%s", wxconf.TokenStorePath, token.wxParams.AppId)
}

func (token *AccessToken) saveToStore() error {
	if _, err := os.Stat(wxconf.TokenStorePath); os.IsNotExist(err) {
		if err = os.MkdirAll(wxconf.TokenStorePath, 0755); err != nil {
			return err
		}
	}
	savePath := token.savePath()
	fp, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	j, _ := json.Marshal(map[string]interface{} {
		"token": token.accessToken,
		"expire": token.expireTime,
	})
	fp.Write(j)
	return nil
}

func (token *AccessToken) loadFromStore() {
	savePath := token.savePath()
	j, err := ioutil.ReadFile(savePath)
	if err != nil {
		return
	}
	var t map[string]interface{}
	if err := json.Unmarshal(j, &t); err != nil {
		return
	}
	if at, ok := t["token"]; ok {
		token.accessToken = at.(string)
	}
	if et, ok := t["expire"]; ok {
		token.expireTime = int64(et.(float64))
	}
}
