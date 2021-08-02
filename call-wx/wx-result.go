package callwx

import (
	"github.com/rosbit/gnet"
	"github.com/rosbit/go-wx-api/v2/log"
	"fmt"
	"net/http"
)

type WxResult interface {
	GetCode() int
	GetMsg() string
}

type BaseResult struct {
	Errcode int
	Errmsg  string
}
func (b *BaseResult) GetCode() int {
	return b.Errcode
}
func (b *BaseResult) GetMsg() string {
	return b.Errmsg
}

type FnCall = gnet.FnCallJ
var (
	HttpCall = gnet.HttpCallJ
	JSONCall = gnet.JSONCallJ
)

func CallWx(url string, method string, params interface{}, headers map[string]string, call FnCall, res WxResult) (code int, err error) {
	status, err := call(url, res, gnet.M(method), gnet.Params(params), gnet.Headers(headers), gnet.BodyLogger(wxlog.GetLogger()))
	if err != nil {
		return -1, err
	}
	if status != http.StatusOK {
		return -2, fmt.Errorf("status %d", status)
	}
	if res.GetCode() != 0 {
		return res.GetCode(), fmt.Errorf("%s", res.GetMsg())
	}
	return 0, nil
}
