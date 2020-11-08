package wxtools

import (
	"github.com/rosbit/go-wx-api/auth"
	"github.com/rosbit/go-wx-api/msg"
	"crypto/sha1"
	"time"
	"fmt"
	"encoding/json"
)

func SignJSAPI(accessToken string, url string) (nonce string, timestamp int64, signature string, err error) {
	u := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)
	var resp []byte
	if resp, err = wxauth.CallWxAPI(u, "GET", nil); err != nil {
		return
	}

	var res struct {
		Errcode  int
		Errmsg   string
		Ticket   string
		ExpiresIn int `json:"expires_in"`
	}
	if err = json.Unmarshal(resp, &res); err != nil {
		return
	}

	if res.Errcode != 0 {
		err = fmt.Errorf("errcode: %d, errmsg: %s", res.Errcode, res.Errmsg)
		return
	}

	nonce = string(wxmsg.GetRandomBytes(16))
	timestamp = time.Now().Unix()

	h := sha1.New()
	fmt.Fprintf(h, "jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", res.Ticket, nonce, timestamp, url)
	signature = fmt.Sprintf("%x", h.Sum(nil))

	return
}
