package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/conf"
	"github.com/rosbit/go-wx-api/v2/auth"
	"fmt"
	"os"
)

func CreateMenu(name string, menuJsonFile string) (error) {
	params := wxconf.GetWxParams(name)
	if params == nil {
		return fmt.Errorf("no params for %s", name)
	}

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
	_, err = wxauth.CallWx(params, genParams, "POST", callwx.JsonCall, &res)
	return err
}

func QueryMenu(name string) (map[string]interface{}, error) {
	return queryMenu(name, "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s")
}

func CurrentSelfmenuInfo(name string) (map[string]interface{}, error) {
	return queryMenu(name, "https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s")
}

func queryMenu(name string, uriFmt string) (map[string]interface{}, error) {
	params := wxconf.GetWxParams(name)
	if params == nil {
		return nil, fmt.Errorf("no params for %s", name)
	}

	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf(uriFmt, accessToken)
		return
	}

	type menu map[string]interface{}
	var res struct {
		callwx.BaseResult
		menu
	}
	if _, err := wxauth.CallWx(params, genParams, "GET", callwx.HttpCall, &res); err != nil {
		return nil, err
	}
	return res.menu, nil
}

func DeleteMenu(name string) (error) {
	params := wxconf.GetWxParams(name)
	if params == nil {
		return fmt.Errorf("no params for %s", name)
	}

	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s", accessToken)
		return
	}

	var res struct {
		callwx.BaseResult
	}
	_, err := wxauth.CallWx(params, genParams, "GET", callwx.HttpCall, &res)
	return err
}

