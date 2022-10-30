package wxmsg

import (
	"github.com/rosbit/go-wx-api/v2/conf"
	"github.com/rosbit/go-wx-api/v2/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var SUCCESS_TEXT = []byte("success")

var MustSignatureArgs = []string{"signature", "timestamp", "nonce"}
const (
	SIGNATURE = iota
	TIMESTAMP
	NONCE
)

type MsgParser interface {
	GetReply(msgBody []byte) (decryptedBody []byte, err error)
	getReply(msgBody []byte) ([]byte, error)
	GetMessageBody(r *http.Request) (msgBody []byte, timestamp string, nonce string, err error)
	getEncryptedMsg(body []byte) (content string, err error)
	EncryptReply(replyMsg []byte, timestamp string, nonce string) []byte
	GetAppId() string
}

type FnGetEnccryptedMsg func(body []byte) (content string, err error)

func parseMessageBody(wxParams *wxconf.WxParamT, u *url.URL, body []byte, getEncryptedMsg FnGetEnccryptedMsg) ([]byte, string, string, error) {
	query := u.Query()
	encrypt_type := query.Get("encrypt_type")
	if encrypt_type == "" {
		return body, "", "", nil
	} else if encrypt_type == "aes" {
		eBody, err := getEncryptedMsg(body)
		if err != nil {
			return nil, "", "", err
		}

		// signautre args are checked in signatrue_checker, so just get them here
		args := make([]string, len(MustSignatureArgs))
		for i, arg := range MustSignatureArgs {
			args[i] = query.Get(arg)
		}

		msg_signature := query.Get("msg_signature")
		msg, err := decryptMsg(wxParams, eBody, msg_signature, args[TIMESTAMP], args[NONCE])
		if err != nil {
			return nil, "", "", err
		}
		wxlog.Logf("plain msg: %s\n", string(msg))
		return msg, args[TIMESTAMP], args[NONCE], nil
	} else {
		return nil, "", "", fmt.Errorf("unsupported encrypted method")
	}
}

type _replyMsg struct {
	reply []byte
	err   error
}
type _reqMsg struct {
	msgBody []byte
	replyChan chan *_replyMsg
}

type msgParserAdapter struct {
	wxParams *wxconf.WxParamT
	msgChan chan *_reqMsg
	mp MsgParser
}

// 根据消息体获取返回消息
func (p *msgParserAdapter) GetReply(msgBody []byte) ([]byte, error) {
	replyChan := make(chan *_replyMsg)
	p.msgChan <- &_reqMsg{msgBody, replyChan}

	replyMsg := <-replyChan
	close(replyChan)

	return replyMsg.reply, replyMsg.err
}

func (p *msgParserAdapter) GetAppId() string {
	return p.wxParams.AppId
}

// 消息解析线程，被GetReply()触发，通过getReply()完成实际的消息处理
func (p *msgParserAdapter) msgParser() {
	for {
		reqMsg := <-p.msgChan
		msgBody, replyChan := reqMsg.msgBody, reqMsg.replyChan

		reply, err := p.mp.getReply(msgBody)
		replyChan <- &_replyMsg{reply, err}
	}
}

// 获取服务号收到的消息参数，返回 (消息体, 时间戳, nonce, error)
func (p *msgParserAdapter) GetMessageBody(r *http.Request) (msgBody []byte, timestapm string, nonce string, err error) {
	if r.Body == nil {
		return nil, "", "", fmt.Errorf("body expected")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, "", "", err
	}
	r.Body.Close()
	wxlog.Logf("body: %s\n", string(body))

	return parseMessageBody(p.wxParams, r.URL, body, p.mp.getEncryptedMsg)
}
