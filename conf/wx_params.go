/**
 * 微信服务号配置信息
 */
package wxconf

import (
	"encoding/base64"
	"fmt"
)

type WxParamT struct {
	Token     string
	AesKey    []byte
	AppId     string
	AppSecret string
}

var (
	wxParams = map[string]*WxParamT{}
	TokenStorePath string
)

func InitTokenStorePath(tokenStorePath string) {
	if len(tokenStorePath) > 0 {
		TokenStorePath = tokenStorePath
	}
}

func (params *WxParamT) setAesKey(aesKey string) error {
	var err error
	if params.AesKey, err = base64.StdEncoding.DecodeString(fmt.Sprintf("%s=", aesKey)); err != nil {
		return err
	}
	if len(params.AesKey) != 32 {
		return fmt.Errorf("invalid wxAESKey length")
	}
	return nil
}

func NewWxParams(name, token, appId, appSecret, aesKey string) error {
	params := &WxParamT{Token:token, AppId:appId, AppSecret:appSecret}
	if aesKey == "" {
		wxParams[name] = params
		return nil
	}
	if err := params.setAesKey(aesKey); err != nil {
		return err
	}
	wxParams[name] = params
	return nil
}

func GetWxParams(name string) (*WxParamT) {
	if params, ok := wxParams[name]; ok {
		return params
	}
	return nil
}

