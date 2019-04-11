/**
 * "subscribe" message handler
 * Rosbit Xu
 */
package main

import (
	"fmt"
	"github.com/rosbit/go-wx-api/msg"
	"github.com/rosbit/go-wx-api/samples/wx-server/utils"
)

// 处理微信用户订阅服务号
func subcribeUser(subscribeEvent *wxmsg.SubscribeEvent) wxmsg.ReplyMsg {
	showWelcome := func(msg string) wxmsg.ReplyMsg {
		if msg != "" {
			return wxmsg.NewReplyTextMsg(subscribeEvent.FromUserName, subscribeEvent.ToUserName, msg)
		}

		if welcome, ok := utils.GetFileContent(WelcomeFile); !ok {
			return wxmsg.NewReplyTextMsg(subscribeEvent.FromUserName, subscribeEvent.ToUserName, "welcome")
		} else {
			return wxmsg.NewReplyTextMsg(subscribeEvent.FromUserName, subscribeEvent.ToUserName, string(welcome))
		}
	}

	res, err := utils.JsonCall(realEventSubscribeUrl, "POST", subscribeEvent)
	if err != nil {
		fmt.Printf("failed to JsonCall(%s): %v\n", realMsgTextUrl, err)
		return showWelcome("")
	}

	if msg, ok := res["msg"]; !ok {
		return showWelcome("")
	} else {
		return showWelcome(msg.(string))
	}
}
