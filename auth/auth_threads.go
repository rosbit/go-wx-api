/**
 * 微信网页授权处理主流程：固定线程数处理所有的请求
 * 工作方式:
 *   1. wxapi.ParseRedirectArgs(*http.Request)用于解析微信授权网页的code和state参数
 *   2. wxapi.AuthRedirect(code, state)用于处理结果
 */
package wxauth

import (
	"net/http"
	"fmt"
	"github.com/rosbit/go-wx-api/conf"
)

type _redirectRes struct {
	msg string
	headers map[string]string
	rurl string
	err error
}

type _redirectData struct {
	code string
	state string
	result chan *_redirectRes
}

type WxAppIdAuthHandler struct {
	wxParams *wxconf.WxParamsT
	reqs chan *_redirectData

	redirectHandler RedirectHandler // 转发处理程序，参见auth_redictor.go
}

func (handler *WxAppIdAuthHandler) GetAppId() string {
	return handler.wxParams.AppId
}

// 微信网页授权处理线程，输入请求被AuthRedirect()触发
func (handler *WxAppIdAuthHandler) authThread() {
	wxUser := NewWxUser(handler.wxParams)
	appId := handler.GetAppId()

	for {
		req := <-handler.reqs
		openId, err := wxUser.GetOpenId(req.code) // 根据请求code获取用户的openId
		if err != nil {
			req.result <- &_redirectRes{"", nil, "", err}
			continue
		}

		redirect := handler.redirectHandler // 线程先启动，handler后设置，所以不能把赋值提到for外面
		msg, headers, rurl, err := redirect(appId, openId, req.state)
		req.result <- &_redirectRes{msg, headers, rurl, err}
	}
}

// 应用初始化时调用，启动若干个线程处理微信网页授权
func StartAuthThreads(params *wxconf.WxParamsT, workerNum int) *WxAppIdAuthHandler {
	h := &WxAppIdAuthHandler{}
	if params == nil {
		h.wxParams = &wxconf.WxParams
	} else {
		h.wxParams = params
	}
	h.reqs = make(chan *_redirectData, workerNum)
	h.redirectHandler = ToAppIdRedirectHandler(HandleRedirect)

	for i:=0; i<workerNum; i++ {
		go h.authThread()
	}
	return h
}

var _mustRedirectArgs = []string{"code", "state"}
const (
	CODE = iota
	STATE
)

// 分析微信网页授权参数，分别返回 (code, state, error)
func ParseRedirectArgs(r *http.Request) (string, string, error) {
	form := r.URL.Query()
	args := make([]string, len(_mustRedirectArgs))
	for i, arg := range _mustRedirectArgs {
		args[i] = form.Get(arg)
		if args[i] == "" {
			return "", "", fmt.Errorf("argument expected")
		}
	}
	return args[CODE], args[STATE], nil
}

// 获取微信网页授权的处理结果，分别返回 [网页内容(非空), 需要设置的header, 跳转的url(非空), error]
func (handler *WxAppIdAuthHandler) AuthRedirect(code string, state string) (string, map[string]string, string, error) {
	rd := &_redirectData{
		code,
		state,
		make(chan *_redirectRes),
	}
	handler.reqs <- rd

	rr := <-rd.result
	close(rd.result)
	return rr.msg, rr.headers, rr.rurl, rr.err
}

