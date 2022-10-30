package wxauth

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/conf"
	"time"
	"fmt"
	"os"
	"encoding/json"
)

type Ticket struct {
	ticket string
	expireTime int64
	wxParams *wxconf.WxParamT
}

func NewTicket(name string) *Ticket {
	wxParams := wxconf.GetWxParams(name)
	if wxParams == nil {
		return nil
	}
	t := &Ticket{
		wxParams: wxParams,
	}
	t.loadFromStore()
	return t
}

func (t *Ticket) Get() (string, error) {
	if t.expired() {
		err := t.get_ticket()
		if err != nil {
			return "", err
		}
	}
	return t.ticket, nil
}

func (t *Ticket) expired() bool {
	return t.expireTime < time.Now().Unix()
}

func (t *Ticket) get_ticket() error {
	token := NewAccessTokenWithParams(t.wxParams)
	accessToken, err := token.Get()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)

	var res struct {
		callwx.BaseResult
		Ticket   string `json:"ticket"`
		ExpiresIn int64 `json:"expires_in"`
	}

	if _, err := callwx.CallWx(url, "GET", nil, nil, callwx.HttpCall, &res); err != nil {
		return err
	}

	t.ticket = res.Ticket
	t.expireTime = res.ExpiresIn + time.Now().Unix() - 10

	return t.saveToStore()
}

func (t *Ticket) savePath() string {
	return fmt.Sprintf("%s/%s.ticket", wxconf.TokenStorePath, t.wxParams.AppId)
}

func (t *Ticket) saveToStore() error {
	if _, err := os.Stat(wxconf.TokenStorePath); os.IsNotExist(err) {
		if err = os.MkdirAll(wxconf.TokenStorePath, 0755); err != nil {
			return err
		}
	}
	savePath := t.savePath()
	fp, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	j, _ := json.Marshal(map[string]interface{} {
		"ticket": t.ticket,
		"expire": t.expireTime,
	})
	fp.Write(j)
	return nil
}

func (t *Ticket) loadFromStore() {
	savePath := t.savePath()
	fp, err := os.Open(savePath)
	if err != nil {
		return
	}
	defer fp.Close()

	var j struct {
		Ticket string `json:"ticket"`
		Expire int64  `json:"expire"`
	}
	dec := json.NewDecoder(fp)
	if err = dec.Decode(&j); err != nil {
		return
	}

	t.ticket, t.expireTime = j.Ticket, j.Expire
}

