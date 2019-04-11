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

func SetAesKey(aesKey string) error {
	var err error
	if WxParams.AesKey, err = base64.StdEncoding.DecodeString(fmt.Sprintf("%s=", aesKey)); err != nil {
		return err
	}
	if len(WxParams.AesKey) != 32 {
		return fmt.Errorf("invalid wxAESKey length")
	}
	return nil
}
