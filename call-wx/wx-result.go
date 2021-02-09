package callwx

import (
	"encoding/json"
	"fmt"
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

func CallWx(url string, method string, params interface{}, headers map[string]string, call FnCallWx, res WxResult) (code int, err error) {
	resp, err := call(url, method, params, headers)
	if err != nil {
		return -1, err
	}
	fmt.Printf("resp: %s\n", resp)

	if err = json.Unmarshal(resp, res); err != nil {
		return -2, err
	}
	if res.GetCode() != 0 {
		return res.GetCode(), fmt.Errorf("%s", res.GetMsg())
	}
	return 0, nil
}
