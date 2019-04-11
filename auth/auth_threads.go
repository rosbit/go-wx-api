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

var (
	reqs chan *_redirectData
)

// 微信网页授权处理线程，输入请求被AuthRedirect()触发
func authThread() {
	wxUser := &WxUser{}

	for {
		req := <-reqs
		openId, err := wxUser.GetOpenId(req.code) // 根据请求code获取用户的openId
		if err != nil {
			req.result <- &_redirectRes{"", nil, "", err}
			continue
		}
		msg, headers, rurl, err := HandleRedirect(openId, req.state) // 调用转发处理程序，参见auth_redictor.go
		req.result <- &_redirectRes{msg, headers, rurl, err}
	}
}

// 应用初始化时调用，启动若干个线程处理微信网页授权
func StartAuthThreads(workerNum int) {
	reqs = make(chan *_redirectData, workerNum)
	for i:=0; i<workerNum; i++ {
		go authThread()
	}
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
func AuthRedirect(code string, state string) (string, map[string]string, string, error) {
	rd := &_redirectData{
		code,
		state,
		make(chan *_redirectRes),
	}
	reqs <- rd

	rr := <-rd.result
	close(rd.result)
	return rr.msg, rr.headers, rr.rurl, rr.err
}

