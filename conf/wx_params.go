/**
 * 微信服务号配置信息
 */
package wxconf

import (
	"encoding/base64"
	"fmt"
)

type WxParamsT struct {
	Token     string
	AesKey    []byte
	AppId     string
	AppSecret string
}

var (
	WxParams WxParamsT
	TokenStorePath string
)

func (params *WxParamsT) SetAesKey(aesKey string) error {
	var err error
	if params.AesKey, err = base64.StdEncoding.DecodeString(fmt.Sprintf("%s=", aesKey)); err != nil {
		return err
	}
	if len(params.AesKey) != 32 {
		return fmt.Errorf("invalid wxAESKey length")
	}
	return nil
}

func NewWxParams(token, appId, appSecret, aesKey string) (*WxParamsT, error) {
	params := &WxParamsT{Token:token, AppId:appId, AppSecret:appSecret}
	if aesKey == "" {
		return params, nil
	}
	if err := params.SetAesKey(aesKey); err != nil {
		return nil, err
	}
	return params, nil
}

func SetAesKey(aesKey string) error {
	return WxParams.SetAesKey(aesKey)
}

func SetParams(token, appId, appSecret, aesKey string) error {
	WxParams.Token, WxParams.AppId, WxParams.AppSecret = token, appId, appSecret
	if aesKey == "" {
		WxParams.AesKey = nil
		return nil
	}
	return WxParams.SetAesKey(aesKey)
}

