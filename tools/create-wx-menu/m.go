package main

import (
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/rosbit/go-wx-api/conf"
	"github.com/rosbit/go-wx-api/auth"
)

type ParamsT struct {
	Token  string  `json:"token"`
	AppId  string
	Secret string
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <params-json-file> <token-cache-path> <menu-json-file>\n", os.Args[0])
		return
	}
	paramsFile, tokenCachePath, menuJsonFile := os.Args[1], os.Args[2], os.Args[3]

	paramsContent, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		fmt.Printf("failed to read params file: %v\n", err)
		return
	}

	var params ParamsT
	if err = json.Unmarshal(paramsContent, &params); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	wxconf.WxParams = wxconf.WxParamsT{Token:params.Token, AppId:params.AppId, AppSecret:params.Secret}
	wxconf.TokenStorePath = tokenCachePath

	menuJson, err := ioutil.ReadFile(menuJsonFile)
	if err != nil {
		fmt.Printf("Failed to read menu json: %v\n", err)
		return
	}

	accessToken, err := wxauth.NewAccessToken().Get()
	if err != nil {
		fmt.Printf("failed to get access token: %v\n", err)
		return
	}

	postUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", accessToken)
	resp, err := wxauth.CallWxAPI(postUrl, "POST", menuJson)
	if err != nil {
		fmt.Printf("failed to call url: %v\n", err)
		return
	}
	fmt.Printf("resp: %s\n", string(resp))
}
