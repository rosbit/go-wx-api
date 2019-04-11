/**
 * 微信消息/事件的缺省处理函数
 * 当前处理的消息/事件未全部覆盖，可以根据需要在此扩充
 */
package wxmsg

var (
	// default message and event handlers

	HandleTextMsg = func (textMsg *TextMsg) ReplyMsg {
		return NewReplyTextMsg(textMsg.FromUserName, textMsg.ToUserName, textMsg.Content)
	}
	HandleViewEvent = func (viewEvent *ViewEvent) ReplyMsg {
		return NewSuccessMsg()
	}
	HandleSubscribeEvent = func (subscribeEvent *SubscribeEvent) ReplyMsg {
		return NewReplyTextMsg(subscribeEvent.FromUserName, subscribeEvent.ToUserName, "welcome")
	}
	HandleUnsubscribeEvent = func (subscribeEvent *SubscribeEvent) ReplyMsg {
		return NewSuccessMsg()
	}
)

