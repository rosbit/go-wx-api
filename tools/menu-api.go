package wxtools

import (
	"github.com/rosbit/go-wx-api/v2/call-wx"
	"github.com/rosbit/go-wx-api/v2/auth"
	"github.com/rosbit/go-wget"
	"fmt"
	"os"
	"encoding/json"
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

func QueryMenu(name string) ([]byte, error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s", accessToken)
		return
	}

	type menu struct {
		Menu            map[string]interface{} `json:"menu"`
		ConditionalMenu map[string]interface{} `json:"conditionalmenu,omitempty"`
	}
	var res struct {
		callwx.BaseResult
		menu
	}
	if _, err := wxauth.CallWx(name, genParams, "GET", wget.HttpCallJ, &res); err != nil {
		return nil, err
	}
	return json.Marshal(res.menu)
}

func CurrentSelfmenuInfo(name string) ([]byte, error) {
	genParams := func(accessToken string)(url string, body interface{}, headers map[string]string) {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s", accessToken)
		return
	}

	type menu struct {
		IsMenuOpen int8 `json:"is_menu_open"`
		SelfMenuInfo map[string]interface{} `json:"selfmenu_info"`
	}
	var res struct {
		callwx.BaseResult
		menu
	}
	if _, err := wxauth.CallWx(name, genParams, "GET", wget.HttpCallJ, &res); err != nil {
		return nil, err
	}
	return json.Marshal(res.menu)
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

