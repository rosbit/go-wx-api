/**
 * 消息/事件的响应消息体构造和序列化为XML，被msg_parse.go调用
 * 目前支持消息类型有限，可以根据需要扩充
 */
package wxmsg

import (
	"fmt"
	"time"
)

// 密文响应体
const _encrypted_reply_form = `<xml>
<Encrypt><![CDATA[%s]]></Encrypt>
    <MsgSignature><![CDATA[%s]]></MsgSignature>
    <TimeStamp>%s</TimeStamp>
    <Nonce><![CDATA[%s]]></Nonce>
</xml>`

// 根据响应消息生成加密消息
func EncryptReply(replyMsg string, timestamp string, nonce string) []byte {
	cryptedText, signature := encryptMsg(replyMsg, timestamp, nonce)
	return []byte(fmt.Sprintf(_encrypted_reply_form, cryptedText, signature, timestamp, nonce))
}

// 响应消息序列化XML接口定义
type ReplyMsg interface {
	ToXML() string
}

type ReplyMsgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
}

// ---------- reply Text Message ---------------
type ReplyTextMsg struct {
	ReplyMsgBase
	Content string
}

func NewReplyTextMsg(toUserName, fromUserName, content string) *ReplyTextMsg {
	return &ReplyTextMsg{ReplyMsgBase:ReplyMsgBase{ToUserName:toUserName, FromUserName:fromUserName, CreateTime:time.Now().Unix()},
		Content:content,
	}
}

func (msg *ReplyTextMsg) ToXML() string {
	return fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>{CreateTime}</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`, msg.ToUserName, msg.FromUserName, msg.Content)
}

// ------------------ reply only "success" ----------
type SuccessTextMsg struct {
}

func NewSuccessMsg() *SuccessTextMsg {
	return &SuccessTextMsg{}
}

func (msg *SuccessTextMsg) ToXML() string {
	return "success"
}
