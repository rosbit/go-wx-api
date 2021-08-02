package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"fmt"
)

func MakeShorturl(name string, longUrl string) (shortUrl string, err error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/shorturl?access_token=%s", accessToken)
		body = map[string]interface{}{
			"action": "long2short",
			"long_url": longUrl,
		}
		return
	}

	var res struct {
		callwx.BaseResult
		ShortUrl string `json:"short_url"`
	}
	if _, err = wxauth.CallWx(name, genParams, "POST", callwx.JSONCall, &res); err != nil {
		return
	}
	shortUrl = res.ShortUrl

	return
}

