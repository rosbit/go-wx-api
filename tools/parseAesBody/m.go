package main

import (
	"github.com/rosbit/go-wx-api/msg"
	"github.com/rosbit/go-wx-api/conf"
	"net/url"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 7 {
		fmt.Printf("Usage: %s <wxToken> <appId> <appSecret> <aesKey> <uri> <xmlBody>\n", os.Args[0])
		return
	}
	token, appId, appSecret, aesKey, uri, xmlBody := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6]

	var err error
	wxconf.WxParams = wxconf.WxParamsT{Token:token, AppId:appId, AppSecret:appSecret}
	if aesKey != "" {
		if err = wxconf.SetAesKey(aesKey); err != nil {
			fmt.Printf("invallid wxAESKey: %v\n", err)
			return
		}
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		fmt.Printf("failed to parse uri: %v\n", err)
		return
	}

	body, timestamp, nonce, err := wxmsg.ParseMessageBody(u, []byte(xmlBody))
	if err != nil {
		fmt.Printf("failed to parse body: %v\n", err)
		return
	}
	fmt.Printf("body: %s\n", string(body))
	fmt.Printf("timestamp: %s\n", timestamp)
	fmt.Printf("nonce: %s\n", nonce)
}
