package wxtools

import (
	"fmt"
	"github.com/rosbit/go-wx-api/auth"
)

func SetTemplateIndustry(accessToken string, industryIds [2]string) ([]byte, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token=%s", accessToken)
	return wxauth.JsonCall(url, "POST", map[string]string{
		"industry_id1": industryIds[0],
		"industry_id2": industryIds[1],
	})
}

func QueryTemplateIndustry(accessToken string) ([]byte, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/template/get_industry?access_token=%s", accessToken)
	return wxauth.CallWxAPI(url, "GET", nil)
}

func SendTemplateMessage(accessToken string, toUser string, templateId string, data map[string]interface{}, url, mpId, mpPagePath string) ([]byte, error) {
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
	u := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)
	return wxauth.JsonCall(u, "POST", d)
}

