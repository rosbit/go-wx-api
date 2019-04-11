/**
 * "text" message handler
 * Rosbit Xu
 */
package main

import (
	"fmt"
	"github.com/rosbit/go-wx-api/msg"
	"github.com/rosbit/go-wx-api/samples/wx-server/utils"
)

// 处理微信用户在服务号中输入的文本消息。如果不处理，返回nil
func textMsgReceived(textMsg *wxmsg.TextMsg) wxmsg.ReplyMsg {
	res, err := utils.JsonCall(realMsgTextUrl, "POST", textMsg)
	if err != nil {
		fmt.Printf("failed to JsonCall(%s): %v\n", realMsgTextUrl, err)
		return nil
	}
	if msg, ok := res["msg"]; !ok {
		fmt.Printf("no \"msg\" item in %v\n", res)
		return nil
	} else {
		return wxmsg.NewReplyTextMsg(textMsg.FromUserName, textMsg.ToUserName, msg.(string))
	}
}
