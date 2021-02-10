/**
 * 微信网页授权处理主流程：固定线程数处理所有的请求
 */
package wxoauth2

import (
	"net/http/httputil"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
	w  http.ResponseWriter
	r *http.Request
	result chan *_redirectRes
}

type WxOAuth2Handler struct {
	service string
	appId string
	reqs chan *_redirectData

	userInfoFlag   string     // redirectUrl是否处理 "snsapi_userinfo" scope的标志串
	redirectUrl    string     // 转发处理URL
}

func (handler *WxOAuth2Handler) redirect(req *_redirectData, wxUser *WxUser) {
	code, state, w, r := req.code, req.state, req.w, req.r

	openId, err := wxUser.GetOpenId(code) // 根据请求code获取用户的openId
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	auth2UserInfo := false
	if len(handler.userInfoFlag) > 0 {
		auth2UserInfo = strings.Index(r.RequestURI, handler.userInfoFlag) >= 0
	}

	var userInfo interface{}
	if auth2UserInfo {
		if err = wxUser.GetInfo(); err == nil {
			userInfo = &wxUser.UserInfo
		}
	} else {
		userInfo, err = wxUser.GetInfoByAccessToken()
	}

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(map[string]interface{}{
		"requestURI": r.RequestURI,
		"appId": handler.appId,
		"openId": openId,
		"state": state,
		"userInfo": userInfo,
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

// 微信网页授权处理线程，输入请求被AuthRedirectUrl()触发
func (handler *WxOAuth2Handler) authThread() {
	wxUser := NewWxUser(handler.service)

	for req := range handler.reqs {
		handler.redirect(req, wxUser)
		req.result <- nil
	}
}

// 应用初始化时调用，启动若干个线程处理微信网页授权
func StartAuthThreads(service, appId string, workerNum int, redirectUrl string, userInfoFlag ...string) *WxOAuth2Handler {
	h := &WxOAuth2Handler{service:service, appId:appId, redirectUrl:redirectUrl}
	h.reqs = make(chan *_redirectData, workerNum)
	if len(userInfoFlag) > 0 {
		h.userInfoFlag = userInfoFlag[0]
	}

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
func (handler *WxOAuth2Handler) AuthRedirectUrl(w http.ResponseWriter, r *http.Request, code, state string) {
	rd := &_redirectData{
		code,
		state,
		w,
		r,
		make(chan *_redirectRes),
	}

	handler.reqs <- rd
	<-rd.result
	close(rd.result)
}
