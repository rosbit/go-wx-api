/**
 * 微信网页授权处理，用于处理服务号菜单点击事件
 * Rosbit Xu
 */
package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/rosbit/go-wx-api/samples/wx-server/utils"
)

var (
	// 缺省的菜单 state -> URL，可以动态更新，根据环境变量 MENU_JSON_CONF 设置
	_state2urls = map[string]string {
		"1": "http://%s/here_are",
		"2": "http://%s/just_some",
		"3": "http://%s/samples.html",
	}

	_menuJsonFileMd5 string
)


func _state2url(state string) (string, bool) {
	f := utils.GetFile(MenuJsonConf)
	if f != nil {
		if _menuJsonFileMd5 != f.Md5sum {
			var m map[string]string
			if err := json.Unmarshal(f.Content, &m); err == nil {
				_state2urls, _menuJsonFileMd5 = m, _menuJsonFileMd5
			}
		}
	}

	rurl, ok := _state2urls[state]
	return rurl, ok
}

/**
 * 根据服务号菜单state做跳转
 * @param openId  订阅用户的openId
 * @param state   微信网页授权中的参数，用来标识某个菜单
 * @return
 *   c    需要显示服务号对话框中的内容
 *   h    需要在微信内嵌浏览器中设置的header信息，包括Cookie
 *   r    需要通过302跳转的URL。如果r不是空串，c的内容被忽略
 *   err  如果没有错误返回nil，非nil表示错误
 */
func handleRedirect(openId, state string) (c string, h map[string]string, r string, err error) {
	if rurl, ok := _state2url(state); !ok {
		err = fmt.Errorf("unknown state %s", state)
		return
	} else {
		if strings.HasPrefix(rurl, "http") {
			// 如果是http(s)打头，表示要跳转的URL
			r = fmt.Sprintf(rurl, MenuHandlerHost)
		} else {
			// 不是URL，表示本地文件，提取内容显示在服务号对话框
			if fc, ok := utils.GetFileContent(rurl); !ok {
				c = rurl
			} else {
				c = string(fc)
			}
			return
		}
	}

	res, e := utils.JsonCall(realMenuRedirectUrl, "POST", map[string]string{"openId": openId, "state": state})
	if e != nil {
		err = e
		return
	}

	if cc, ok := res["c"]; ok {
		c = cc.(string)
	}
	if hh, ok := res["h"]; ok {
		h1 := hh.(map[string]interface{})
		h = make(map[string]string, len(h1))
		for k, v := range h1 {
			h[k] = v.(string)
		}
	}
	if rr, ok := res["r"]; ok {
		r = rr.(string)
	}

	return
}
