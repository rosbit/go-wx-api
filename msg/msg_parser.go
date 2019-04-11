/**
 * 服务号消息/事件的结构定义和解析，触发消息/事件处理
 * 1. wxapi.ParseMessageBody()  --- 获取消息/事件的(消息体,时间戳,nonce,error)，主要用于调试
 * 2. wxapi.GetReply([]byte)    --- 根据消息体触发消息处理函数得到返回结果。如果使用default_rest.go的实现不需要关心它
 */
package wxmsg

import (
	"github.com/beevik/etree"
	"fmt"
	"strconv"
	"net/http"
	"io/ioutil"
	"net/url"
	"github.com/rosbit/go-wx-api/log"
)

const (
	SUCCESS_TEXT = "success"

	// msgType
	MT_TEXT  = "text"
	MT_EVENT = "event"

	// event type
	ET_SUBSCRIBE   = "subscribe"
	ET_UNSUBSCRIBE = "unsubscribe"
	ET_VIEW        = "VIEW"
)

var MustSignatureArgs = []string{"signature", "timestamp", "nonce"}
const (
	SIGNATURE = iota
	TIMESTAMP
	NONCE
)

func _getText(el *etree.Element, tagName string) (string, error) {
	t := el.SelectElement(tagName)
	if t == nil {
		return "", fmt.Errorf("%s not found", tagName)
	}
	return t.Text(), nil
}

// --------------------- message -------------------
type BaseMsg struct {
	ToUserName string
	FromUserName string
	CreateTime int
	MsgType string
}

func (m *BaseMsg) parse(root *etree.Element) {
	m.ToUserName, _   = _getText(root, "ToUserName")
	m.FromUserName, _ = _getText(root, "FromUserName")
	m.MsgType, _      = _getText(root, "MsgType")
	ct, _            := _getText(root, "CreateTime")
	m.CreateTime, _ = strconv.Atoi(ct)
}

type Msg struct {
	BaseMsg
	MsgId string
}

func (m *Msg) parse(root *etree.Element) {
	m.BaseMsg.parse(root)
	m.MsgId, _ = _getText(root, "MsgId")
}

type TextMsg struct {
	Msg
	Content string
}

func (m *TextMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.Content, _ = _getText(root, "Content")
}

// ------------- event ---------------
type EventMsg struct {
	Msg
	Event string
}

func (m *EventMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.Event, _ = _getText(root, "Event")
}

// VIEW
type ViewEvent struct {
	EventMsg
	EventKey string
	MenuId string
}

func (m *ViewEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.EventKey, _ = _getText(root, "EventKey")
	m.MenuId, _   = _getText(root, "MenuId")
}

// subscribe, unsubscribe
type SubscribeEvent struct {
	EventMsg
	EventKey string
	Ticket string
}

func (m *SubscribeEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.EventKey, _ = _getText(root, "EventKey")
	m.Ticket, _   = _getText(root, "Ticket")
}

// 消息/事件主处理流程：分析消息内容、根据消息类型触发消息处理函数、返回结果消息
func getReply(msgBody []byte) (string, error) {
	msg := etree.NewDocument()
	err := msg.ReadFromBytes(msgBody)
	if err != nil {
		return SUCCESS_TEXT, err
	}

	root := msg.SelectElement("xml")
	msgType, _ := _getText(root, "MsgType")
	var replyMsg ReplyMsg

	switch msgType {
	case MT_TEXT:
		textMsg := &TextMsg{}
		textMsg.parse(root)
		replyMsg = HandleTextMsg(textMsg)
	case MT_EVENT:
		eventType, _ := _getText(root, "Event")
		switch eventType {
		case ET_VIEW:
			viewEvent := &ViewEvent{}
			viewEvent.parse(root)
			replyMsg = HandleViewEvent(viewEvent)
		case ET_SUBSCRIBE:
			subscribeEvent := &SubscribeEvent{}
			subscribeEvent.parse(root)
			replyMsg = HandleSubscribeEvent(subscribeEvent)
		case ET_UNSUBSCRIBE:
			subscribeEvent := &SubscribeEvent{}
			subscribeEvent.parse(root)
			replyMsg = HandleUnsubscribeEvent(subscribeEvent)
		default:
			return SUCCESS_TEXT, fmt.Errorf("under implementation for event type: %s", eventType)
		}
	default:
		return SUCCESS_TEXT, fmt.Errorf("under implementation for msg type: %s", msgType)
	}
	if replyMsg == nil {
		return SUCCESS_TEXT, nil
	}
	return replyMsg.ToXML(), nil
}

type _replyMsg struct {
	reply string
	err   error
}
type _reqMsg struct {
	msgBody []byte
	replyChan chan *_replyMsg
}

var _msgChan chan *_reqMsg

// 消息解析线程，被GetReply()触发，通过getReply()完成实际的消息处理
func msgParser() {
	for {
		reqMsg := <-_msgChan
		msgBody, replyChan := reqMsg.msgBody, reqMsg.replyChan

		reply, err := getReply(msgBody)
		replyChan <- &_replyMsg{reply, err}
	}
}

// 初始化应用时启动若干个消息解析线程
func StartWxMsgParsers(workNum int) {
	_msgChan = make(chan *_reqMsg, workNum)
	for i:=0; i<workNum; i++ {
		go msgParser()
	}
}

// 根据消息体获取返回消息
func GetReply(msgBody []byte) (string, error) {
	replyChan := make(chan *_replyMsg)
	_msgChan <- &_reqMsg{msgBody, replyChan}

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
func ParseMessageBody(u *url.URL, body []byte) ([]byte, string, string, error) {
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
		msg, err := decryptMsg(eBody, msg_signature, args[TIMESTAMP], args[NONCE])
		if err != nil {
			return nil, "", "", err
		}
		return msg, args[TIMESTAMP], args[NONCE], nil
	} else {
		return nil, "", "", fmt.Errorf("unsupported encrypted method")
	}
}

// 获取服务号收到的消息参数，返回 (消息体, 时间戳, nonce, error)
func GetMessageBody(r *http.Request) ([]byte, string, string, error) {
	if r.Body == nil {
		return nil, "", "", fmt.Errorf("body expected")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, "", "", err
	}
	r.Body.Close()
	wxlog.Logf("body: %s\n", string(body))

	return ParseMessageBody(r.URL, body)
}
