package callwx

import (
	"github.com/rosbit/go-wget"
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

func CallWx(url string, method string, params interface{}, headers map[string]string, call wget.FnCallJ, res WxResult) (code int, err error) {
	status, err := call(url, method, params, headers, res)
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
