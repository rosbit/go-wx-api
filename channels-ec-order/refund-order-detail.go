// 获取售后单
// 参考文档: https://developers.weixin.qq.com/doc/channels/API/aftersale/getaftersaleorder.html

package ceord

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"fmt"
	"encoding/json"
)

func GetRefundOrderDetail(appName string, afterSaleOrderId string) (ord json.RawMessage, err error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/channels/ec/aftersale/getaftersaleorder?access_token=%s", accessToken)
		body = map[string]interface{}{
			"after_sale_order_id": afterSaleOrderId,
		}
		return
	}

	var res struct {
		callwx.BaseResult
		AfterSaleOrder json.RawMessage `json:"after_sale_order"`
	}

	if _, err = wxauth.CallWx(appName, genParams, "POST", callwx.JSONCall, &res, true); err != nil {
		return
	}
	ord = res.AfterSaleOrder
	return
}

