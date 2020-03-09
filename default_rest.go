/**
 * RESTful API implementation
 *  GET  /wx         -- 用于实现服务号接口配置，实际路径通过http路由关联
 *  POST /wx         -- 处理微信消息/事件的入口，实际路径通过http路由关联
 *  GET  /redirect   -- 微信网页授权接口，处理服务号菜单的总入口，
 *                      不同菜单通过网页授权参数state区分，实际路径通过http路由关联
 * Rosbit Xu
 */
package wxapi

import (
	"net/http"
	"github.com/rosbit/go-wx-api/log"
	"github.com/rosbit/go-wx-api/auth"
)

// 用于微信服务号设置
func Echo(w http.ResponseWriter, r *http.Request) {
	defaultWxHandler.Echo(w, r)
}

// 微信服务号消息/事件处理入口
func Request(w http.ResponseWriter, r *http.Request) {
	defaultWxHandler.Request(w, r)
}

// 微信网页授权
func Redirect(w http.ResponseWriter, r *http.Request) {
	defaultWxHandler.Redirect(w, r)
}

func _WriteMessage(w http.ResponseWriter, msg string) {
	w.Write([]byte(msg))
}

func _WriteBytes(w http.ResponseWriter, msg []byte) {
	w.Write(msg)
}

func (wx *WxHandler) Echo(w http.ResponseWriter, r *http.Request) {
	wxlog.Logf("wxapi.Echo for appId %s called: %s\n", wx.appIdHandler.GetAppId(), r.RequestURI)
	q := r.URL.Query()
	echostr := q.Get("echostr")
	if echostr != "" {
		_WriteMessage(w, echostr)
	} else {
		_WriteMessage(w, "hello, this is handler view")
	}
}

func (wx *WxHandler) Request(w http.ResponseWriter, r *http.Request) {
	wxlog.Logf("wxapi.Request for appId %s called: %s\n", wx.appIdHandler.GetAppId(), r.RequestURI)
	msgBody, timestamp, nonce, err := wx.appMsgParser.GetMessageBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if msgBody == nil {
		http.Error(w, "no message", http.StatusBadRequest)
		return
	}

	replyMsg, err := wx.appMsgParser.GetReply(msgBody)
	if err != nil || nonce == "" {
		_WriteBytes(w, replyMsg)
		return
	}
	_WriteBytes(w, wx.appMsgParser.EncryptReply(replyMsg, timestamp, nonce))
}

func (wx *WxHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	wxlog.Logf("wxapi.Redirect for appId %s called: %s\n", wx.appIdHandler.GetAppId(), r.RequestURI)
	code, state, err := wxauth.ParseRedirectArgs(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if wx.appIdHandler.HasRedirectUrl() {
		wx.appIdHandler.AuthRedirectUrl(w, r, code, state)
		return
	}

	msg, headers, rurl, err := wx.appIdHandler.AuthRedirect(code, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if headers != nil {
		for k,v := range headers {
			w.Header().Add(k, v)
		}
	}
	if rurl != "" {
		http.Redirect(w, r, rurl, http.StatusFound)
		return
	}
	_WriteMessage(w, msg)
}
