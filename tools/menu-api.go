package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"github.com/rosbit/go-wget"
	"fmt"
	"os"
)

func CreateMenu(name string, menuJsonFile string) (error) {
	fpMenuJson, err := os.Open(menuJsonFile)
	if err != nil {
		return err
	}
	defer fpMenuJson.Close()

	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", accessToken)
		body = fpMenuJson
		return
	}

	var res struct {
		callwx.BaseResult
	}
	_, err = wxauth.CallWx(name, genParams, "POST", wget.JsonCallJ, &res)
	return err
}

func QueryMenu(name string) (map[string]interface{}, error) {
	return queryMenu(name, "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s")
}

func CurrentSelfmenuInfo(name string) (map[string]interface{}, error) {
	return queryMenu(name, "https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s")
}

func queryMenu(name string, uriFmt string) (map[string]interface{}, error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf(uriFmt, accessToken)
		return
	}

	type menu map[string]interface{}
	var res struct {
		callwx.BaseResult
		menu
	}
	if _, err := wxauth.CallWx(name, genParams, "GET", wget.HttpCallJ, &res); err != nil {
		return nil, err
	}
	return res.menu, nil
}

func DeleteMenu(name string) (error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s", accessToken)
		return
	}

	var res struct {
		callwx.BaseResult
	}
	_, err := wxauth.CallWx(name, genParams, "GET", wget.HttpCallJ, &res)
	return err
}

