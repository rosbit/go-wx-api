// 所有需要access_token调用的统一入口

package wxauth

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"fmt"
)

type FnGeneParams func(accessToken string)(url string, body interface{}, headers map[string]string)

func CallWx(name string, genParams FnGeneParams, method string, call callwx.FnCall, res callwx.WxResult, checkChannelsEc ...bool) (code int, err error) {
	token := NewAccessToken(name)
	if token == nil {
		err = fmt.Errorf("no params for %s", name)
		return
	}
	if len(checkChannelsEc) > 0 && checkChannelsEc[0] {
		if !token.wxParams.IsChannelsEc {
			err = fmt.Errorf("params %s is not for channels ec", name)
			return
		}
	}
	accessToken, err := token.Get()
	if err != nil {
		return 0, err
	}
	url, body, headers := genParams(accessToken)
	return callwx.CallWx(url, method, body, headers, call, res)
}
