package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"github.com/rosbit/go-wget"
	"fmt"
)

func SetTemplateIndustry(name string, industryIds [2]string) (map[string]interface{}, error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token=%s", accessToken)
		body = map[string]interface{}{
			"industry_id1": industryIds[0],
			"industry_id2": industryIds[1],
		}
		return
	}
	return templateAction(name, genParams, "POST", wget.JsonCallJ)
}

func QueryTemplateIndustry(name string) (map[string]interface{}, error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/get_industry?access_token=%s", accessToken)
		return
	}
	return templateAction(name, genParams, "GET", wget.HttpCallJ)
}

func SendTemplateMessage(name string, toUser string, templateId string, data map[string]interface{}, url, mpId, mpPagePath string) (map[string]interface{}, error) {
	dData := make(map[string]interface{})
	for k,v := range data {
		dData[k] = map[string]string{"value": fmt.Sprintf("%v", v)}
	}

	d := map[string]interface{}{
		"touser": toUser,
		"template_id": templateId,
		"data": dData,
	}
	if len(url) > 0 {
		d["url"] = url
	}
	if len(mpId) > 0 && len(mpPagePath) > 0 {
		d["miniprogram"] = map[string]string{
			"appid": mpId,
			"pagepath": mpPagePath,
		}
	}

	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)
		body = d
		return
	}
	return templateAction(name, genParams, "POST", wget.JsonCallJ)
}

func templateAction(name string, genParams wxauth.FnGeneParams, method string, call wget.FnCallJ) (map[string]interface{}, error) {
	type result map[string]interface{}
	var res struct {
		callwx.BaseResult
		result
	}
	if _, err := wxauth.CallWx(name, genParams, method, call, &res); err != nil {
		return nil, err
	}
	return res.result, nil
}
