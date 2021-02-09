package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"github.com/rosbit/go-wx-api/v2/conf"
	"fmt"
)

func CreateTempQrIntScene(name string, sceneId int, expireInSec int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createTempQr(name, sceneId, expireInSec, "QR_SCENE", "scene_id")
	return
}

func CreateTempQrStrScene(name, sceneId string, expireInSec int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createTempQr(name, sceneId, expireInSec, "QR_STR_SCENE", "scene_str")
	return
}

func createTempQr(name string, sceneId interface{}, expireInSec int, action, idName string) (string, string, error) {
	params := map[string]interface{}{
		"expire_seconds": expireInSec,
		"action_name": action,
		"action_info": map[string]interface{}{
			"scene": map[string]interface{}{
				idName : sceneId,
			},
		},
	}
	return createQr(name, params)
}

func CreateQrIntScene(name string, sceneId int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createForeverQr(name, sceneId, "QR_LIMIT_SCENE", "scene_id")
	return
}

func CreateQrStrScene(name, sceneId string) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createForeverQr(name, sceneId, "QR_LIMIT_STR_SCENE", "scene_str")
	return
}

func createForeverQr(name string, sceneId interface{}, action, idName string) (string, string, error) {
	params := map[string]interface{}{
		"action_name": action,
		"action_info": map[string]interface{}{
			"scene": map[string]interface{}{
				idName: sceneId,
			},
		},
	}
	return createQr(name, params)
}

func createQr(name string, params map[string]interface{}) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	wxParams := wxconf.GetWxParams(name)
	if wxParams == nil {
		err = fmt.Errorf("no params for %s", name)
		return
	}

	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", accessToken)
		body = params
		return
	}

	var res struct {
		callwx.BaseResult
		Ticket        string `json:"ticket"`
		ExpireSeconds int    `json:"expire_seconds"`
		Url           string `json:"url"`
	}
	if _, err = wxauth.CallWx(wxParams, genParams, "POST", callwx.JsonCall, &res); err != nil {
		return
	}

	if res.Ticket == "" {
		err = fmt.Errorf("no ticket item found in result")
		return
	}
	ticketURL2ShowQrCode = fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", res.Ticket)

	if res.Url == "" {
		err = fmt.Errorf("no url item found in result")
		return
	}
	urlIncluedInQrcode = res.Url
	return
}
