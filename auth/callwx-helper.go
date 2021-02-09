// 所有需要access_token调用的统一入口

package wxauth

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/conf"
)

type FnGeneParams func(accessToken string)(url string, body interface{}, headers map[string]string)

func CallWx(wxParams *wxconf.WxParamT, genParams FnGeneParams, method string, call callwx.FnCallWx, res callwx.WxResult) (code int, err error) {
	accessToken, err := NewAccessToken(wxParams).Get()
	if err != nil {
		return 0, err
	}
	url, body, headers := genParams(accessToken)
	return callwx.CallWx(url, method, body, headers, call, res)
}
