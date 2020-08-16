package wxtools

import (
	"github.com/rosbit/go-wx-api/auth"
	"fmt"
	"encoding/json"
)

func MakeShorturl(accessToken string, longUrl string) (shortUrl string, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/shorturl?access_token=%s", accessToken)
	var b []byte
	if b, err = wxauth.JsonCall(url, "POST", map[string]string{
		"action": "long2short",
		"long_url": longUrl,
	}); err != nil {
		return
	}

	var res struct {
		Errcode int
		Errmsg string
		ShortUrl string `json:"short_url"`
	}
	if err = json.Unmarshal(b, &res); err != nil {
		return
	}
	if res.Errcode != 0 {
		err = fmt.Errorf("errcode: %d, errmsg: %s", res.Errcode, res.Errmsg)
		return
	}
	shortUrl = res.ShortUrl
	return
}

