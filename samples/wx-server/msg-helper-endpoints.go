/**
 * real service endpoint
 */
package main

import (
	"fmt"
	"strings"
)

const (
	REAL_MSG_TEXT          = "/msg/text"
	REAL_EVENT_SUBSCRIBE   = "/event/subscribe"
	REAL_EVENT_UNSUBSCRIBE = "/event/unsubscribe"
	MENU_REDIRECT          = "/menu/redirect"
)

var (
	realMsgTextUrl          string
	realEventSubscribeUrl   string
	realEventUnsubscribeUrl string
	realMenuRedirectUrl     string
)

func _toUrl(prefix, endpoint string) string {
	url := fmt.Sprintf("%s%s", prefix, endpoint)
	fmt.Printf(" + %s\n", url)
	return url
}

func createMsgHelperEndpoints() {
	var prefix string
	if strings.HasSuffix(MsgHelperPrefix, "/") {
		prefix = MsgHelperPrefix[:len(MsgHelperPrefix)-1]
	} else {
		prefix = MsgHelperPrefix
	}

	fmt.Printf("message helper endpoints:\n")
	realMsgTextUrl          = _toUrl(prefix, REAL_MSG_TEXT)
	realEventSubscribeUrl   = _toUrl(prefix, REAL_EVENT_SUBSCRIBE)
	realEventUnsubscribeUrl = _toUrl(prefix, REAL_EVENT_UNSUBSCRIBE)
	realMenuRedirectUrl     = _toUrl(prefix, MENU_REDIRECT)
}
