package callwx

import (
	"github.com/rosbit/go-wget"
	"fmt"
)

type FnCallWx func(url string, method string, params interface{}, headers map[string]string) (resp []byte, err error)

func HttpCall(url string, method string, postData interface{}, headers map[string]string) ([]byte, error) {
	return callWget(url, method, postData, headers, wget.Wget)
}

func JsonCall(url string, method string, jsonData interface{}, headers map[string]string) ([]byte, error) {
	return callWget(url, method, jsonData, headers, wget.PostJson)
}

func callWget(url string, method string, postData interface{}, headers map[string]string, fnCall wget.HttpFunc) ([]byte, error) {
	status, content, _, err := fnCall(url, method, postData, headers)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("status %d", status)
	}
	return content, nil
}
