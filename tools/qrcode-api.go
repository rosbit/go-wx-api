package wxtools

import (
	"fmt"
	"encoding/json"
	"github.com/rosbit/go-wx-api/auth"
)

func CreateTempQRIntScene(accessToken string, sceneId int, expireInSec int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createTempQr(accessToken, sceneId, expireInSec, "QR_SCENE", "scene_id")
	return
}

func CreateTempQRStrScene(accessToken, sceneId string, expireInSec int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
	ticketURL2ShowQrCode, urlIncluedInQrcode, err = createTempQr(accessToken, sceneId, expireInSec, "QR_LIMIT_STR_SCENE", "scene_str")
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

func CreateQRIntScene(accessToken string, sceneId int) (ticketURL2ShowQrCode, urlIncluedInQrcode string, err error) {
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
	var postData []byte
	var resp []byte
	if postData, err = json.Marshal(params); err != nil {
		return
	}
	if resp, err = wxauth.CallWxAPI(url, "POST", postData); err != nil {
		return
	}
	var res map[string]interface{}
	if err = json.Unmarshal(resp, &res); err != nil {
		return
	}
	if ticket, ok := res["ticket"]; !ok {
		err = fmt.Errorf("no ticket item found in result")
		return
	} else {
		ticketURL2ShowQrCode = fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", ticket.(string))
	}
	if urlInQrcode, ok := res["url"]; !ok {
		err = fmt.Errorf("no url item found in result")
		return
	} else {
		urlIncluedInQrcode = urlInQrcode.(string)
	}
	return
}
