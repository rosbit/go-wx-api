/**
 * 触发消息/事件处理
 * 1. wxapi.ParseMessageBody()  --- 获取消息/事件的(消息体,时间戳,nonce,error)，主要用于调试
 * 2. wxapi.GetReply([]byte)    --- 根据消息体触发消息处理函数得到返回结果。如果使用default_rest.go的实现不需要关心它
 */
package wxmsg

import (
	"github.com/beevik/etree"
	"fmt"
	"net/http"
	"io/ioutil"
	"net/url"
	"github.com/rosbit/go-wx-api/log"
	"github.com/rosbit/go-wx-api/conf"
)

var SUCCESS_TEXT = []byte("success")

var MustSignatureArgs = []string{"signature", "timestamp", "nonce"}
const (
	SIGNATURE = iota
	TIMESTAMP
	NONCE
)

// msgType => ReceivedMsg
var newMessages = map[string]func()ReceivedMsg {
	MT_TEXT:       func()ReceivedMsg { return &TextMsg{} },
	MT_IMAGE:      func()ReceivedMsg { return &ImageMsg{} },
	MT_VOICE:      func()ReceivedMsg { return &VoiceMsg{} },
	MT_VIDEO:      func()ReceivedMsg { return &VideoMsg{} },
	MT_SHORTVIDEO: func()ReceivedMsg { return &VideoMsg{} },
	MT_LOCATION:   func()ReceivedMsg { return &LocationMsg{} },
	MT_LINK:       func()ReceivedMsg { return &LinkMsg{} },
}

// eventType => ReceivedMsg
var newEvents = map[string]func()ReceivedMsg {
	ET_VIEW:       func()ReceivedMsg { return &ViewEvent{} },
	ET_CLICK:      func()ReceivedMsg { return &ClickEvent{} },
	ET_SUBSCRIBE:  func()ReceivedMsg { return &SubscribeEvent{} },
	ET_UNSUBSCRIBE:func()ReceivedMsg { return &SubscribeEvent{} },
	ET_SCAN:       func()ReceivedMsg { return &SubscribeEvent{} },
	ET_WHERE:      func()ReceivedMsg { return &WhereEvent{} },
	ET_LOCATION:   func()ReceivedMsg { return &LocationEvent{} },
	ET_PIC_SYSPHOTO:       func()ReceivedMsg { return &PhotoEvent{} },
	ET_PIC_PHOTO_OR_ALBUM: func()ReceivedMsg { return &PhotoEvent{} },
	ET_PIC_WEIXIN:         func()ReceivedMsg { return &PhotoEvent{} },
	ET_SCANCODE_WAITMSG:   func()ReceivedMsg { return &ScancodeEvent{} },
	ET_SCANCODE_PUSH:      func()ReceivedMsg { return &ScancodeEvent{} },
	ET_MASSSENDJOBFINISH:     func()ReceivedMsg { return &MassSentEvent{} },
	ET_TEMPLATESENDJOBFINISH: func()ReceivedMsg { return &TemplateSentEvent{} },
}

// 消息/事件主处理流程：分析消息内容、根据消息类型触发消息处理函数、返回结果消息
func (p *WxAppIdMsgParser) getReply(msgBody []byte) ([]byte, error) {
	msg := etree.NewDocument()
	err := msg.ReadFromBytes(msgBody)
	if err != nil {
		return SUCCESS_TEXT, err
	}

	root := msg.SelectElement("xml")
	msgType, _ := _getText(root, "MsgType")
	var eventType string

	var replyMsg ReplyMsg
	var receivedMsg ReceivedMsg
	if msgType != MT_EVENT {
		if newMessge, ok := newMessages[msgType]; ok {
			receivedMsg = newMessge()
			receivedMsg.parse(root)
			replyMsg = p.handleReceivedMessage(receivedMsg, msgType)
		} else {
			return SUCCESS_TEXT, fmt.Errorf("under implementation for msg type: %s", msgType)
		}
	} else {
		eventType, _ = _getText(root, "Event")
		if newEvent, ok := newEvents[eventType]; ok {
			receivedMsg = newEvent()
			receivedMsg.parse(root)
			replyMsg = p.handleReceivedEvent(receivedMsg, eventType)
		} else {
			return SUCCESS_TEXT, fmt.Errorf("under implementation for event type: %s", eventType)
		}
	}

	if replyMsg == nil {
		return SUCCESS_TEXT, nil
	}
	return replyMsg.ToXML(), nil
}

type _replyMsg struct {
	reply []byte
	err   error
}
type _reqMsg struct {
	msgBody []byte
	replyChan chan *_replyMsg
}

type WxAppIdMsgParser struct {
	wxParams *wxconf.WxParamsT
	msgChan chan *_reqMsg

	messageHandlers map[string]FnMessageHandler
	eventHandlers   map[string]FnMessageHandler
}

// 消息解析线程，被GetReply()触发，通过getReply()完成实际的消息处理
func (p *WxAppIdMsgParser) msgParser() {
	for {
		reqMsg := <-p.msgChan
		msgBody, replyChan := reqMsg.msgBody, reqMsg.replyChan

		reply, err := p.getReply(msgBody)
		replyChan <- &_replyMsg{reply, err}
	}
}

// 初始化应用时启动若干个消息解析线程
func StartWxMsgParsers(params *wxconf.WxParamsT, workNum int) *WxAppIdMsgParser {
	p := &WxAppIdMsgParser{}
	if params == nil {
		p.wxParams = &wxconf.WxParams
	} else {
		p.wxParams = params
	}
	p.RegisterWxMsgHandler(MsgHandler) // set default msg handler.

	p.msgChan = make(chan *_reqMsg, workNum)
	for i:=0; i<workNum; i++ {
		go p.msgParser()
	}

	return p
}

// 根据消息体获取返回消息
func (p *WxAppIdMsgParser) GetReply(msgBody []byte) ([]byte, error) {
	replyChan := make(chan *_replyMsg)
	p.msgChan <- &_reqMsg{msgBody, replyChan}

	replyMsg := <-replyChan
	close(replyChan)

	return replyMsg.reply, replyMsg.err
}

func getEncryptedMsg(body []byte) (string, error) {
	msg := etree.NewDocument()
	if err := msg.ReadFromBytes(body); err != nil {
		return "", err
	}

	root := msg.SelectElement("xml")
	return _getText(root, "Encrypt")
}

// 从GetMessageBody()独立出来，可以通过各种方式调用，方便调试
func parseMessageBody(wxParams *wxconf.WxParamsT, u *url.URL, body []byte) ([]byte, string, string, error) {
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

// 使用缺省配置解析消息，可以用于调试
func ParseMessageBody(u *url.URL, body []byte) ([]byte, string, string, error) {
	return parseMessageBody(&wxconf.WxParams, u, body)
}

// 获取服务号收到的消息参数，返回 (消息体, 时间戳, nonce, error)
func (p *WxAppIdMsgParser) GetMessageBody(r *http.Request) ([]byte, string, string, error) {
	if r.Body == nil {
		return nil, "", "", fmt.Errorf("body expected")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, "", "", err
	}
	r.Body.Close()
	wxlog.Logf("body: %s\n", string(body))

	return parseMessageBody(p.wxParams, r.URL, body)
}
