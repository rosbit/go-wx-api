// 获取订单详情
// 参考文档: https://developers.weixin.qq.com/doc/channels/API/order/get.html

package ceord

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"fmt"
	"encoding/json"
)

func GetOrderDetail(appName string, orderId string) (ord json.RawMessage, err error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/channels/ec/order/get?access_token=%s", accessToken)
		body = map[string]interface{}{
			"order_id": orderId,
		}
		return
	}

	var res struct {
		callwx.BaseResult
		Order json.RawMessage `json:"order"`
	}

	if _, err = wxauth.CallWx(appName, genParams, "POST", callwx.JSONCall, &res, true); err != nil {
		return
	}
	ord = res.Order
	return
}

