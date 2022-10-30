/**
 * RESTful API implementation
 * Rosbit Xu
 */
package wxapi

import (
	"github.com/rosbit/go-wx-api/v2/oauth2"
	"github.com/rosbit/go-wx-api/v2/msg"
	"github.com/rosbit/go-wx-api/v2/conf"
	"github.com/rosbit/go-wx-api/v2/log"
	"fmt"
	"net/http"
)

// 用于微信服务号设置
// 路由方法: GET
// uri?signature=xxx=timestamp=xxx&nonce=xxx&echostr=xxx
func CreateEcho(token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		signature := q.Get("signature")
		timestamp := q.Get("timestamp")
		nonce := q.Get("nonce")
		echostr := q.Get("echostr")
		var h string
		if len(signature) == 0 || len(timestamp) == 0 || len(nonce) == 0 || len(echostr) == 0 {
			goto ERR
		}

		h = wxmsg.HashStrings([]string{
			token, timestamp, nonce,
		})
		if signature != h {
			goto ERR
		}

		fmt.Fprintf(w, "%s", echostr)
		return
ERR:
		fmt.Fprintf(w, "hello, this is handler view")
	}
}

// 创建微信服务号消息/事件处理入口
// 路由方法: POST
// @param serviceName  配置项的名称
// @parma workerNum    处理消息的并发数
// @param msgHandler   消息/事件处理器，根据实际情况实现
func CreateMsgHandler(serviceName string, workerNum int, msgHandler wxmsg.WxMsgHandler) http.HandlerFunc {
	params := wxconf.GetWxParams(serviceName)
	if params == nil {
		panic(fmt.Errorf("params for %s not found", serviceName))
	}
	appMsgParser := wxmsg.StartWxMsgParsers(params, workerNum)
	appMsgParser.RegisterWxMsgHandler(msgHandler)

	return func(w http.ResponseWriter, r *http.Request) {
		wxlog.Logf("wxapi.Request for appId %s called: %s\n", params.AppId, r.RequestURI)
		handleMsg(w, r, appMsgParser)
	}
}

// 创建视频号小店事件处理入口
// 路由方法: POST
// @param serviceName  配置项的名称
// @parma workerNum    处理消息的并发数
// @param channelsEcEventHandler 视频号小店事件处理器，根据实际情况实现
func CreateChannelsEcHandler(serviceName string, workerNum int, channelsEcEventHandler wxmsg.ChannelsEcEventHandler) http.HandlerFunc {
	params := wxconf.GetWxParams(serviceName)
	if params == nil {
		panic(fmt.Errorf("params for %s not found", serviceName))
	}

	channelsEcEventParser := wxmsg.StartChannelsEcParsers(params, workerNum)
	channelsEcEventParser.RegisterChannelsEcEventHandler(channelsEcEventHandler)

	return func(w http.ResponseWriter, r *http.Request) {
		wxlog.Logf("wxapi.Request for appId %s called: %s\n", params.AppId, r.RequestURI)
		handleMsg(w, r, channelsEcEventParser)
	}
}

func handleMsg(w http.ResponseWriter, r *http.Request, msgParser wxmsg.MsgParser) {
	msgBody, timestamp, nonce, err := msgParser.GetMessageBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if msgBody == nil {
		http.Error(w, "no message", http.StatusBadRequest)
		return
	}

	replyMsg, err := msgParser.GetReply(msgBody)
	if err != nil || nonce == "" {
		w.Write(replyMsg)
		return
	}
	w.Write(msgParser.EncryptReply(replyMsg, timestamp, nonce))
}

// 创建网页授权处理器
// 路由方法: GET
// @param serviceName  配置项的名称
// @parma workerNum    处理消息的并发数
// @param redirectUrl  该URL将全权决定网页授权的处理
//     请求方式: POST
//     请求BODY: 是一个JSON: {"appId": "xxx", "openId": "xxx", "state": "state"}
//     该URL的以POST形式接收参数，而且会得到所有的HTTP头信息，可以设置任何的响应头信息
//     响应结果直接显示在公众号浏览器中，响应时间要控制好，避免微信服务超时
// @param userInfoFlag 只取第一项，用于检查转发url中是否有标志串；该值存在表示使用 snsapi_userinfo 获取用户信息
func CreateOAuth2Redirector(serviceName string, workerNum int, redirectUrl string, userInfoFlag ...string) http.HandlerFunc {
	params := wxconf.GetWxParams(serviceName)
	if params == nil {
		panic(fmt.Errorf("params for %s not found", serviceName))
	}
	oauth2Redirector := wxoauth2.StartAuthThreads(serviceName, params.AppId, workerNum, redirectUrl, userInfoFlag...)

	return func(w http.ResponseWriter, r *http.Request) {
		wxlog.Logf("wxapi.Redirect for appId %s called: %s\n", params.AppId, r.RequestURI)
		code, state, err := wxoauth2.ParseRedirectArgs(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		oauth2Redirector.AuthRedirectUrl(w, r, code, state)
	}
}
