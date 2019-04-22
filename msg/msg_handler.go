/**
 * 微信消息/事件的缺省处理函数
 */
package wxmsg

import (
	"fmt"
)

// 消息/事件处理接口定义
type WxMsgHandler interface {
	HandleTextMsg(textMsg *TextMsg) ReplyMsg
	HandleImageMsg(imgMsg *ImageMsg) ReplyMsg
	HandleVoiceMsg(voiceMsg *VoiceMsg) ReplyMsg
	HandleVideoMsg(videoMsg *VideoMsg) ReplyMsg
	HandleLocationMsg(locMsg *LocationMsg) ReplyMsg
	HandleLinkMsg(linkMsg *LinkMsg) ReplyMsg
	HandleClickEvent(clickEvent *ClickEvent) ReplyMsg
	HandleViewEvent(viewEvent *ViewEvent) ReplyMsg
	HandleScanEvent(scanEvent *ScanEvent) ReplyMsg
	HandleScanWaitEvent(scanEvent *ScanEvent)ReplyMsg
	HandleSubscribeEvent(subscribeEvent *SubscribeEvent) ReplyMsg
	HandleUnsubscribeEvent(subscribeEvent *SubscribeEvent) ReplyMsg
	HandleWhereEvent(whereEvent *WhereEvent) ReplyMsg
	HandlePhotoEvent(phoneEvent *PhotoEvent) ReplyMsg
	HandleLocationEvent(locEvent *LocationEvent) ReplyMsg
	HandleMassSentEvent(massSentEvent *MassSentEvent) ReplyMsg
	HandleTemplateSentEvent(tempSentEvent *TemplateSentEvent) ReplyMsg
}

// 缺省的消息/事件处理器，可以通过RegisterWxMsghandler覆盖
var MsgHandler WxMsgHandler = &WxMsgHandlerAdapter{}

// 缺省的消息处理器实现。如果要重新实现某些处理方法，需要
// 1. 在新的结构体中嵌入(embed)该结构体
// 2. 覆盖实现某些方法
// 3. 调用RegisterWxMsghandler覆盖缺省实现
type WxMsgHandlerAdapter struct {
}

func (h *WxMsgHandlerAdapter) HandleTextMsg(textMsg *TextMsg) ReplyMsg {
	return NewReplyTextMsg(textMsg.FromUserName, textMsg.ToUserName, textMsg.Content)
}

func (h *WxMsgHandlerAdapter) HandleImageMsg(imgMsg *ImageMsg) ReplyMsg {
	return NewReplyImageMsg(imgMsg.FromUserName, imgMsg.ToUserName, imgMsg.MediaId)
}

func (h *WxMsgHandlerAdapter) HandleVoiceMsg(voiceMsg *VoiceMsg) ReplyMsg {
	return NewReplyVoiceMsg(voiceMsg.FromUserName, voiceMsg.ToUserName, voiceMsg.MediaId)
}

func (h *WxMsgHandlerAdapter) HandleVideoMsg(videoMsg *VideoMsg) ReplyMsg {
	return NewReplyVideoMsg(videoMsg.FromUserName, videoMsg.ToUserName, videoMsg.MediaId, "echo", "thank you")
}

func (h *WxMsgHandlerAdapter) HandleLocationMsg(locMsg *LocationMsg) ReplyMsg {
	return NewSuccessMsg()
}

func (h *WxMsgHandlerAdapter) HandleLinkMsg(linkMsg *LinkMsg) ReplyMsg {
	return NewReplyTextMsg(linkMsg.FromUserName, linkMsg.ToUserName, fmt.Sprintf("your link %s received", linkMsg.Title))
}

func (h *WxMsgHandlerAdapter) HandleClickEvent(clickEvent *ClickEvent) ReplyMsg {
	return NewReplyTextMsg(clickEvent.FromUserName, clickEvent.ToUserName,  fmt.Sprintf("event key %s clicked", clickEvent.EventKey))
}

func (h *WxMsgHandlerAdapter) HandleViewEvent(viewEvent *ViewEvent) ReplyMsg {
	return NewSuccessMsg()
}

func (h *WxMsgHandlerAdapter) HandleScanEvent(scanEvent *ScanEvent) ReplyMsg {
	return NewSuccessMsg()
}

func (h *WxMsgHandlerAdapter) HandleScanWaitEvent(scanEvent *ScanEvent) ReplyMsg {
	return NewReplyTextMsg(scanEvent.FromUserName, scanEvent.ToUserName, fmt.Sprintf("scan result: %s, %s", scanEvent.ScanType, scanEvent.ScanResult))
}

func (h *WxMsgHandlerAdapter) HandleSubscribeEvent(subscribeEvent *SubscribeEvent) ReplyMsg {
	return NewReplyTextMsg(subscribeEvent.FromUserName, subscribeEvent.ToUserName, "welcome")
}

func (h *WxMsgHandlerAdapter) HandleUnsubscribeEvent(subscribeEvent *SubscribeEvent) ReplyMsg {
	return NewSuccessMsg()
}

func (h *WxMsgHandlerAdapter) HandleWhereEvent(whereEvent *WhereEvent) ReplyMsg {
	return NewSuccessMsg()
}

func (h *WxMsgHandlerAdapter) HandlePhotoEvent(phoneEvent *PhotoEvent) ReplyMsg {
	return NewSuccessMsg()
}

func (h *WxMsgHandlerAdapter) HandleLocationEvent(locEvent *LocationEvent) ReplyMsg {
	return NewSuccessMsg()
}

func (h *WxMsgHandlerAdapter) HandleMassSentEvent(massSentEvent *MassSentEvent) ReplyMsg {
	return NewReplyTextMsg(massSentEvent.FromUserName, massSentEvent.ToUserName, massSentEvent.Status)
}

func (h *WxMsgHandlerAdapter) HandleTemplateSentEvent(tempSentEvent *TemplateSentEvent) ReplyMsg {
	return NewReplyTextMsg(tempSentEvent.FromUserName, tempSentEvent.ToUserName, tempSentEvent.Status)
}

type FnMessageHandler func(receivedMsg ReceivedMsg) ReplyMsg

func (p *WxAppIdMsgParser) handleReceivedMessage(receivedMsg ReceivedMsg, msgType string) ReplyMsg {
	if fn, ok := p.messageHandlers[msgType]; ok {
		return fn(receivedMsg)
	}
	return nil
}

func (p *WxAppIdMsgParser) handleReceivedEvent(receivedMsg ReceivedMsg, eventType string) ReplyMsg {
	if fn, ok := p.eventHandlers[eventType]; ok {
		return fn(receivedMsg)
	}
	return nil
}

func (p *WxAppIdMsgParser) RegisterWxMsgHandler(msgHandler WxMsgHandler) {
	if msgHandler == nil {
		return
	}

	p.messageHandlers = map[string]FnMessageHandler {
		MT_TEXT:       func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleTextMsg(receivedMsg.(*TextMsg)) },
		MT_IMAGE:      func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleImageMsg(receivedMsg.(*ImageMsg)) },
		MT_VOICE:      func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleVoiceMsg(receivedMsg.(*VoiceMsg)) },
		MT_SHORTVIDEO: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleVideoMsg(receivedMsg.(*VideoMsg)) },
		MT_LOCATION:   func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleLocationMsg(receivedMsg.(*LocationMsg)) },
		MT_LINK:       func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleLinkMsg(receivedMsg.(*LinkMsg)) },
	}

	p.eventHandlers = map[string]FnMessageHandler {
		ET_VIEW:  func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleViewEvent(receivedMsg.(*ViewEvent)) },
		ET_CLICK: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleClickEvent(receivedMsg.(*ClickEvent)) },
		ET_SUBSCRIBE:   func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleSubscribeEvent(receivedMsg.(*SubscribeEvent)) },
		ET_UNSUBSCRIBE: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleUnsubscribeEvent(receivedMsg.(*SubscribeEvent)) },
		ET_WHERE:    func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleWhereEvent(receivedMsg.(*WhereEvent)) },
		ET_LOCATION: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleLocationEvent(receivedMsg.(*LocationEvent)) },
		ET_PIC_SYSPHOTO: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandlePhotoEvent(receivedMsg.(*PhotoEvent)) },
		ET_PIC_PHOTO_OR_ALBUM: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandlePhotoEvent(receivedMsg.(*PhotoEvent)) },
		ET_PIC_WEIXIN: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandlePhotoEvent(receivedMsg.(*PhotoEvent)) },
		ET_SCANCODE_WAITMSG: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleScanWaitEvent(receivedMsg.(*ScanEvent)) },
		ET_SCANCODE_PUSH: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleScanEvent(receivedMsg.(*ScanEvent)) },
		ET_MASSSENDJOBFINISH:func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleMassSentEvent(receivedMsg.(*MassSentEvent)) },
		ET_TEMPLATESENDJOBFINISH: func(receivedMsg ReceivedMsg) ReplyMsg { return msgHandler.HandleTemplateSentEvent(receivedMsg.(*TemplateSentEvent)) },
	}
}

