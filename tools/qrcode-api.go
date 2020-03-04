package wxtools

import (
	"fmt"
	"encoding/json"
	"github.com/rosbit/go-wx-api/auth"
)

func CreateTempQrIntScene(accessToken string, sceneId int, expireInSec int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createTempQr(accessToken, sceneId, expireInSec, "QR_SCENE", "scene_id")
	return
}

func CreateTempQrStrScene(accessToken, sceneId string, expireInSec int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createTempQr(accessToken, sceneId, expireInSec, "QR_STR_SCENE", "scene_str")
	return
}

func createTempQr(accessToken string, sceneId interface{}, expireInSec int, action, idName string) (string, string, error) {
	params := map[string]interface{}{
		"expire_seconds": expireInSec,
		"action_name": action,
		"action_info": map[string]interface{}{
			"scene": map[string]interface{}{
				idName : sceneId,
			},
		},
	}
	return createQr(accessToken, params)
}

func CreateQrIntScene(accessToken string, sceneId int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createForeverQr(accessToken, sceneId, "QR_LIMIT_SCENE", "scene_id")
	return
}

func CreateQrStrScene(accessToken, sceneId string) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createForeverQr(accessToken, sceneId, "QR_LIMIT_STR_SCENE", "scene_str")
	return
}

func createForeverQr(accessToken string, sceneId interface{}, action, idName string) (string, string, error) {
	params := map[string]interface{}{
		"action_name": action,
		"action_info": map[string]interface{}{
			"scene": map[string]interface{}{
				idName: sceneId,
			},
		},
	}
	return createQr(accessToken, params)
}

func createQr(accessToken string, params map[string]interface{}) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", accessToken)
	var resp []byte
	if resp, err = wxauth.JsonCall(url, "POST", params); err != nil {
		return
	}

	var res struct {
		Ticket        string `json:"ticket"`
		ExpireSeconds int    `json:"expire_seconds"`
		Url           string `json:"url"`
	}
	if err = json.Unmarshal(resp, &res); err != nil {
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
