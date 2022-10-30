/**
 * 触发消息/事件处理
 *  wxapi.GetReply([]byte)    --- 根据消息体触发消息处理函数得到返回结果。如果使用default_rest.go的实现不需要关心它
 */
package wxmsg

import (
	"github.com/rosbit/go-wx-api/v2/conf"
	"github.com/beevik/etree"
	"fmt"
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

type WxAppIdMsgParser struct {
	*msgParserAdapter

	messageHandlers map[string]FnMessageHandler
	eventHandlers   map[string]FnMessageHandler
}

// 初始化应用时启动若干个消息解析线程
func StartWxMsgParsers(params *wxconf.WxParamT, workNum int) *WxAppIdMsgParser {
	p := &WxAppIdMsgParser{
		msgParserAdapter: &msgParserAdapter{
			wxParams:params,
		},
	}
	p.msgParserAdapter.mp = p
	p.RegisterWxMsgHandler(MsgHandler) // set default msg handler.

	p.msgChan = make(chan *_reqMsg, workNum)
	for i:=0; i<workNum; i++ {
		go p.msgParser()
	}

	return p
}

func (p *WxAppIdMsgParser) getEncryptedMsg(body []byte) (string, error) {
	msg := etree.NewDocument()
	if err := msg.ReadFromBytes(body); err != nil {
		return "", err
	}

	root := msg.SelectElement("xml")
	return _getText(root, "Encrypt")
}
