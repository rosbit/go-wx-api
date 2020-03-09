/**
 * 微信网页授权处理主流程：固定线程数处理所有的请求
 * 工作方式:
 *   1. wxapi.ParseRedirectArgs(*http.Request)用于解析微信授权网页的code和state参数
 *   2. wxapi.AuthRedirect(code, state)用于处理结果
 *   3. wxapi.AuthRedirectUrl(http.ResponseWriter, *http.Request, code, state)用于全权处理网页授权，优先级高于AuthRediret()
 */
package wxauth

import (
	"github.com/rosbit/go-wx-api/conf"
	"net/http/httputil"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
	"fmt"
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

	redirectUrl     string          // 转发处理URL，**如果存在，下面的redirectHandler将被忽略**
	redirectHandler RedirectHandler // 转发处理程序，参见auth_redictor.go
}

func (handler *WxAppIdAuthHandler) GetAppId() string {
	return handler.wxParams.AppId
}

func (handler *WxAppIdAuthHandler) HasRedirectUrl() bool {
	return len(handler.redirectUrl) > 0
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

// 网页授权全权转发给redirectUrl
func (handler *WxAppIdAuthHandler) AuthRedirectUrl(w http.ResponseWriter, r *http.Request, code, state string) {
	wxUser := NewWxUser(handler.wxParams)
	openId, err := wxUser.GetOpenId(code) // 根据请求code获取用户的openId
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = wxUser.GetInfo()

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(map[string]interface{}{
		"appId": handler.wxParams.AppId,
		"openId": openId,
		"state": state,
		"userInfo": &(wxUser.UserInfo),
		"userInfoError": func(err error)string{
			if err == nil {
				return ""
			} else {
				return err.Error()
			}
		}(err),
	})

	forwarder := func()*httputil.ReverseProxy{
		return &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.Method = http.MethodPost
				r.URL, _ = url.Parse(handler.redirectUrl)
				r.Body = ioutil.NopCloser(b)
				r.ContentLength = int64(b.Len())
				r.Header.Set("Content-Type", "application/json")
			},
		}
	}()

	forwarder.ServeHTTP(w, r)
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

