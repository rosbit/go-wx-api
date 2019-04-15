package wxtools

import (
	"fmt"
	"net/url"
	"io/ioutil"
	"github.com/rosbit/go-wx-api/auth"
)

func CreateMenu(accessToken string, menuJsonFile string) ([]byte, error) {
	menuJson, err := ioutil.ReadFile(menuJsonFile)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", accessToken)
	return wxauth.CallWxAPI(url, "POST", menuJson)
}

func QueryMenu(accessToken string) ([]byte, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s", accessToken)
	return wxauth.CallWxAPI(url, "GET", nil)
}

func DeleteMenu(accessToken string) ([]byte, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s", accessToken)
	return wxauth.CallWxAPI(url, "GET", nil)
}

func CurrentSelfmenuInfo(accessToken string) ([]byte, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s", accessToken)
	return wxauth.CallWxAPI(url, "GET", nil)
}

var _valid_scopes = map[string]bool{"snsapi_base":true, "snsapi_userinfo":true}

func isValidScope(scope string) bool {
	_, ok := _valid_scopes[scope]
	return ok
}

func MakeAuthUrl(appId, redirectUri, scope, state string) (string, error) {
	if !isValidScope(scope) {
		return "", fmt.Errorf("unknown scope value: %s", scope)
	}
	encUri := url.QueryEscape(redirectUri)
	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect", appId, encUri, scope, state), nil
}
