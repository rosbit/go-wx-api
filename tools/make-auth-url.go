package wxtools

import (
	"net/url"
	"fmt"
)

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
