package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/auth"
	"github.com/rosbit/go-wx-api/v2/msg"
	"crypto/sha1"
	"time"
	"fmt"
)

func SignJSAPI(name string, url string) (nonce string, timestamp int64, signature string, err error) {
	t := wxauth.NewTicket(name)
	if t == nil {
		err = fmt.Errorf("no conf found for %s", name)
		return
	}
	var ticket string
	if ticket, err = t.Get(); err != nil {
		return
	}

	nonce = string(wxmsg.GetRandomBytes(16))
	timestamp = time.Now().Unix()

	h := sha1.New()
	fmt.Fprintf(h, "jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, nonce, timestamp, url)
	signature = fmt.Sprintf("%x", h.Sum(nil))

	return
}
