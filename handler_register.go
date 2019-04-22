/**
 * 注册微信消息/事件处理器，用于覆盖缺省处理器
 */
package wxapi

import (
	"github.com/rosbit/go-wx-api/msg"
	"github.com/rosbit/go-wx-api/auth"
)

// 函数签名定义，各消息结构定义见msg_receive.go，返回结果结构定义见msg_reply.go

/**
 * 注册消息/事件处理器
 * @msgHandler  消息处理器
 */
func RegisterWxMsghandler(msgHandler wxmsg.WxMsgHandler) {
	defaultWxHandler.RegisterWxMsghandler(msgHandler)
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
type RedirectHandler = wxauth.RedirectHandlerWithoutAppId
//type RedirectHandler func(openId, state string) (c string, h map[string]string, r string, err error)

/**
 * 注册微信网页授权处理函数
 */
func RegisterRedictHandler(handler RedirectHandler) {
	if handler != nil {
		defaultWxHandler.RegisterRedictHandler(wxauth.ToAppIdRedirectHandler(handler))
	}
}

// ---------------- 支持多服务号 ------------------
func (h *WxHandler) RegisterWxMsghandler(msgHandler wxmsg.WxMsgHandler) {
	h.appMsgParser.RegisterWxMsgHandler(msgHandler)
}

func (h *WxHandler) RegisterRedictHandler(handler wxauth.RedirectHandler) {
	h.appIdHandler.RegisterRedictHandler(handler)
}
