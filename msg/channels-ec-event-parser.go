/**
 * 触发视频号小店事件处理
 *  wxapi.GetReply([]byte)    --- 根据消息体触发消息处理函数得到返回结果。如果使用default_rest.go的实现不需要关心它
 */
package wxmsg

import (
	"github.com/rosbit/go-wx-api/v2/conf"
	"fmt"
	"encoding/json"
)

// eventType => ReceivedJSONEvent
var newChannelsEcEvents = map[string]func()ReceivedJSONEvent {
	ET_ORDER_CANCEL: func()ReceivedJSONEvent { return &OrderCancelEvent{} },
	ET_ORDER_PAY:    func()ReceivedJSONEvent { return &OrderPayEvent{} },
	ET_ORER_CONFIRM: func()ReceivedJSONEvent { return &OrderConfirmEvent{} },
	ET_ORDER_SETTLE: func()ReceivedJSONEvent { return &OrderSettleEvent{} },
	ET_AFTERSAL_UPDATE: func()ReceivedJSONEvent { return &AftersaleUpdateEvent{} },
}

// 视频号小店事件主处理流程：分析消息内容、根据事件类型触发事件处理函数、返回结果消息
func (p *ChannelsEcEventParser) getReply(msgBody []byte) ([]byte, error) {
	var e ChannelsEcEvent
	if err := json.Unmarshal(msgBody, &e); err != nil {
		return SUCCESS_TEXT, err
	}

	msgType := e.MsgType
	var eventType string

	var replyMsg []byte
	var receivedMsg ReceivedJSONEvent
	if msgType != MT_EVENT {
		return SUCCESS_TEXT, fmt.Errorf("under implementation for msg type: %s", msgType)
	} else {
		eventType = e.Event
		if newEvent, ok := newChannelsEcEvents[eventType]; ok {
			receivedMsg = newEvent()
			if err := json.Unmarshal(msgBody, receivedMsg); err != nil {
				return SUCCESS_TEXT, err
			}
			replyMsg = p.handleReceivedEvent(receivedMsg, eventType)
		} else {
			return SUCCESS_TEXT, fmt.Errorf("under implementation for event type: %s", eventType)
		}
	}

	if replyMsg == nil {
		return SUCCESS_TEXT, nil
	}
	return replyMsg, nil
}

type ChannelsEcEventParser struct {
	*msgParserAdapter

	eventHandlers   map[string]FnJSONEventHandler
}

// 初始化应用时启动若干个消息解析线程
func StartChannelsEcParsers(params *wxconf.WxParamT, workNum int) *ChannelsEcEventParser {
	p := &ChannelsEcEventParser{
		msgParserAdapter: &msgParserAdapter{
			wxParams:params,
		},
	}
	p.msgParserAdapter.mp = p
	p.RegisterChannelsEcEventHandler(CEEventHandler) // set default msg handler.

	p.msgChan = make(chan *_reqMsg, workNum)
	for i:=0; i<workNum; i++ {
		go p.msgParser()
	}

	return p
}

func (p *ChannelsEcEventParser) getEncryptedMsg(body []byte) (string, error) {
	var res struct {
		Encrypt string
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	return res.Encrypt, nil
}

func (p *ChannelsEcEventParser) EncryptReply(replyMsg []byte, timestamp string, nonce string) []byte {
	return replyMsg
}
