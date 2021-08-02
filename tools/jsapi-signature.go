package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"github.com/rosbit/go-wx-api/v2/msg"
	"crypto/sha1"
	"time"
	"fmt"
)

func SignJSAPI(name string, url string) (nonce string, timestamp int64, signature string, err error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)
		return
	}

	var res struct {
		callwx.BaseResult
		Ticket   string
		ExpiresIn int `json:"expires_in"`
	}
	if _, err = wxauth.CallWx(name, genParams, "GET", callwx.HttpCall, &res); err != nil {
		return
	}

	nonce = string(wxmsg.GetRandomBytes(16))
	timestamp = time.Now().Unix()

	h := sha1.New()
	fmt.Fprintf(h, "jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", res.Ticket, nonce, timestamp, url)
	signature = fmt.Sprintf("%x", h.Sum(nil))

	return
}
