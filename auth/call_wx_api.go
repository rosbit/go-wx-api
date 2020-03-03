package wxauth

import (
	"github.com/rosbit/go-wget"
	"fmt"
)

func CallWxAPI(url string, method string, postData interface{}) ([]byte, error) {
	return callWget(url, method, postData, wget.Wget)
}

func JsonCall(url string, method string, jsonData interface{}) ([]byte, error) {
	return callWget(url, method, jsonData, wget.PostJson)
}

func callWget(url string, method string, postData interface{}, fnCall wget.HttpFunc) ([]byte, error) {
	status, content, _, err := fnCall(url, method, postData, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("status %d", status)
	}
	return content, nil
}
