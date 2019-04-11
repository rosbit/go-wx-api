package wxauth

import (
	"github.com/rosbit/go-wget"
	"fmt"
)

func CallWxAPI(url string, method string, postData interface{}) ([]byte, error) {
	status, content, _, err := wget.Wget(url, method, postData, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("status %d", status)
	}
	return content, nil
}
