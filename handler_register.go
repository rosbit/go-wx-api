/**
 * 注册微信消息/事件处理函数，用于覆盖缺省处理函数
 * 当前处理的消息/事件未全部覆盖，可以根据需要在此扩充
 */
package wxapi

import (
	"github.com/rosbit/go-wx-api/msg"
	"github.com/rosbit/go-wx-api/auth"
)

// 函数签名定义，各消息结构定义见msg_parser.go，返回结果结构定义见msg_reply.go

type TextMsgHandler func(*wxmsg.TextMsg) wxmsg.ReplyMsg                     // text消息处理
type ViewEventHandler func(*wxmsg.ViewEvent) wxmsg.ReplyMsg                 // VIEW事件处理
type SubscribeEventHandler func(*wxmsg.SubscribeEvent) wxmsg.ReplyMsg       // 用户关注(subscribe)服务号事件处理
type UnsubscribeEventHandler func(*wxmsg.SubscribeEvent) wxmsg.ReplyMsg     // 用户取消关注(unsubscribe)服务号事件处理

// 注册消息/事件处理函数

func RegisterTextMsgHandler(handler TextMsgHandler) {
	wxmsg.HandleTextMsg = handler
}

func RegisterViewEventHandler(handler ViewEventHandler) {
	wxmsg.HandleViewEvent = handler
}

func RegisterSubscribeEventHandler(handler SubscribeEventHandler) {
	wxmsg.HandleSubscribeEvent = handler
}

func UnregisterSubscribeEventHandler(handler UnsubscribeEventHandler) {
	wxmsg.HandleUnsubscribeEvent = handler
}

/**
 * [函数签名]根据服务号菜单state做跳转
 * @param openId  订阅用户的openId
 * @param state   微信网页授权中的参数，用来标识某个菜单
 * @return
 *   c    需要显示服务号对话框中的内容
 *   h    需要在微信内嵌浏览器中设置的header信息，包括Cookie
 *   r    需要通过302跳转的URL。如果r不是空串，c的内容被忽略
 *   err  如果没有错误返回nil，非nil表示错误
 */
type RedirectHandler func(openId, state string) (c string, h map[string]string, r string, err error)

/**
 * 注册微信网页授权处理函数
 */
func RegisterRedictHandler(handler RedirectHandler) {
	wxauth.HandleRedirect = handler
}
