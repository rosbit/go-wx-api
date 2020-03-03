/**
 * signature checker as a http middleware
 * Rosbit Xu
 */
package wxapi

import (
	"net/http"
	"strconv"
	"time"
	"strings"
	"github.com/rosbit/go-wx-api/msg"
)

/**
 * 创建http处理中间件，验证消息签名，如果非法直接返回错误
 * @param wxToken      公众号在微信管理后台定义的token
 * @param timeout      消息时间戳超时处理，秒数，如果<=0不检查时间戳
 * @param uriPrefixes  需要检查签名的URI前缀列表，不相关的URI忽略检查；如果为nil，全部检查
 */
func NewWxSignatureChecker(wxToken string, timeout int, uriPrefixes []string) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if uriPrefixes != nil {
			uri := r.URL.Path
			found := false
			for _, prefix := range uriPrefixes {
				if strings.HasPrefix(uri, prefix) {
					found = true
					break
				}
			}
			if !found {
				next(w, r)
				return
			}
		}
		args := make([]string, len(wxmsg.MustSignatureArgs))
		query := r.URL.Query()
		for i, arg := range wxmsg.MustSignatureArgs {
			args[i] = query.Get(arg)
			if args[i] == "" {
				http.Error(w, "argument expected", http.StatusBadRequest)
				return
			}
		}

		if timeout > 0 {
			ts, err := strconv.ParseInt(args[wxmsg.TIMESTAMP], 10, 64)
			if err != nil {
				http.Error(w, "invalid timestamp", http.StatusBadRequest)
				return
			}
			if ts + int64(timeout) < time.Now().Unix() {
				http.Error(w, "signature expired", http.StatusBadRequest)
				return
			}
		}
		l := []string{wxToken, args[wxmsg.TIMESTAMP], args[wxmsg.NONCE]}

		hashcode := wxmsg.HashStrings(l)
		if hashcode != args[wxmsg.SIGNATURE] {
			http.Error(w, "invalid signanure", http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}
