/**
 * 消息/事件的响应消息体构造和序列化为XML，被msg_parse.go调用
 */
package wxmsg

import (
	"fmt"
	"time"
	"bytes"
	"text/template"
)

const (
	// reply template name
	RTN_ENCRYPTED_REPLY = iota
	RTN_TEXT
	RTN_IMAGE
	RTN_VOICE
	RTN_VIDEO
	RTN_MUSIC
	RTN_NEWS
	RTN_TOTAL
)

var strTmpls = map[int]string {
	RTN_ENCRYPTED_REPLY: `<xml>
<Encrypt><![CDATA[{{.encryptedText}}]]></Encrypt>
    <MsgSignature><![CDATA[{{.signature}}]]></MsgSignature>
    <TimeStamp>{{.timestamp}}</TimeStamp>
    <Nonce><![CDATA[{{.nonce}}]]></Nonce>
</xml>`,
	RTN_TEXT: `<xml>
<ToUserName><![CDATA[{{.ToUserName}}]]></ToUserName>
<FromUserName><![CDATA[{{.FromUserName}}]]></FromUserName>
<CreateTime>{{.CreateTime}}</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[{{.Content}}]]></Content>
</xml>`,
	RTN_IMAGE: `<xml>
<ToUserName><![CDATA[{{.ToUserName}}]]></ToUserName>
<FromUserName><![CDATA[{{.FromUserName}}]]></FromUserName>
<CreateTime>{{.CreateTime}}</CreateTime>
<MsgType><![CDATA[image]]></MsgType>
<Image>
<MediaId><![CDATA[{{.MediaId}}]]></MediaId>
</Image>
</xml>`,
	RTN_VOICE: `<xml>
<ToUserName><![CDATA[{{.ToUserName}}]]></ToUserName>
<FromUserName><![CDATA[{{.FromUserName}}]]></FromUserName>
<CreateTime>{{.CreateTime}}</CreateTime>
<MsgType><![CDATA[voice]]></MsgType>
<Voice>
<MediaId><![CDATA[{{.MediaId}}]]></MediaId>
</Voice>
</xml>`,
	RTN_VIDEO: `<xml>
<ToUserName><![CDATA[{{.ToUserName}}]]></ToUserName>
<FromUserName><![CDATA[{{.FromUserName}}]]></FromUserName>
<CreateTime>{{.CreateTime}}</CreateTime>
<MsgType><![CDATA[video]]></MsgType>
<Video>
<MediaId><![CDATA[{{.MediaId}}]]></MediaId>
<Title><![CDATA[{{.Title}}]]></Title>
<Description><![CDATA[{{.Desc}}]]></Description>
</Video>
</xml>`,
	RTN_MUSIC: `<xml>
<ToUserName><![CDATA[{{.ToUserName}}]]></ToUserName>
<FromUserName><![CDATA[{{.FromUserName}}]]></FromUserName>
<CreateTime>{{.CreateTime}}</CreateTime>
<MsgType><![CDATA[music]]></MsgType><Music>
<Title><![CDATA[{{.Title}}]]></Title>
<Description><![CDATA[{{.Desc}}]]></Description>
<MusicUrl><![CDATA[{{.MusicUrl}}]]></MusicUrl>
<HQMusicUrl><![CDATA[{{.HQMusicUrl}}]]></HQMusicUrl>
<ThumbMediaId><![CDATA[{{.ThumbMediaId}}]]></ThumbMediaId>
</Music>
</xml>`,
	RTN_NEWS: `<xml>
<ToUserName><![CDATA[{{.ToUserName}}]]></ToUserName>
<FromUserName><![CDATA[{{.FromUserName}}]]></FromUserName>
<CreateTime>{{.CreateTime}}</CreateTime>
<MsgType><![CDATA[news]]></MsgType>
<ArticleCount>{{ len .Articles }}</ArticleCount>
<Articles>
{{range .Articles}}
	<item>
		<Title><![CDATA[{{.Title}}]]></Title>
		<Description><![CDATA[{{.Desc}}]]></Description>
		<PicUrl><![CDATA[{{.PicUrl}}]]></PicUrl>
		<Url><![CDATA[{{.Url}}]]></Url>
	</item>
{{end}}
</Articles>
</xml>`,
}

var replyTmpls []*template.Template

func init() {
	replyTmpls = make([]*template.Template, RTN_TOTAL)
	var err error
	for idx, tmpl := range strTmpls {
		replyTmpls[idx], err = template.New("ROSBIT").Parse(tmpl)
		if err != nil {
			fmt.Printf("failed to parse template #%d, %s: %v\n", idx, tmpl, err)
		}
	}
}

func executeTmpl(tmplIdx int, data interface{}) []byte {
	bb := bytes.Buffer{}
	replyTmpls[tmplIdx].Execute(&bb, data)
	return bb.Bytes()
}

// 根据响应消息生成加密消息
func (p *WxAppIdMsgParser) EncryptReply(replyMsg []byte, timestamp string, nonce string) []byte {
	cryptedText, signature := encryptMsg(p.wxParams, replyMsg, timestamp, nonce)
	return executeTmpl(RTN_ENCRYPTED_REPLY, map[string]string{
		"encryptedText": cryptedText,
		"signature": signature,
		"timestamp": timestamp,
		"nonce": nonce,
	})
}

// 响应消息序列化XML接口定义
type ReplyMsg interface {
	ToXML() []byte
}

type ReplyMsgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
}

func newMsgBase(toUserName, fromUserName string) ReplyMsgBase {
	return ReplyMsgBase{ToUserName:toUserName, FromUserName:fromUserName, CreateTime:time.Now().Unix()}
}

// ---------- reply Text Message ---------------
type ReplyTextMsg struct {
	ReplyMsgBase
	Content string
}

func NewReplyTextMsg(toUserName, fromUserName, content string) *ReplyTextMsg {
	return &ReplyTextMsg{
		ReplyMsgBase:newMsgBase(toUserName, fromUserName),
		Content:content,
	}
}

func (msg *ReplyTextMsg) ToXML() []byte {
	return executeTmpl(RTN_TEXT, msg)
}

// ------------------ reply only "success" ----------
type SuccessTextMsg struct {
}

func NewSuccessMsg() *SuccessTextMsg {
	return &SuccessTextMsg{}
}

func (msg *SuccessTextMsg) ToXML() []byte {
	return SUCCESS_TEXT
}

// ----------------- reply image message ----------
type ReplyImageMsg struct {
	ReplyMsgBase
	MediaId string
}

func NewReplyImageMsg(toUserName, fromUserName, mediaId string) *ReplyImageMsg {
	return &ReplyImageMsg{
		ReplyMsgBase:newMsgBase(toUserName, fromUserName),
		MediaId:mediaId,
	}
}

func (msg *ReplyImageMsg) ToXML() []byte {
	return executeTmpl(RTN_IMAGE, msg)
}

// ----------------- reply void message ----------------
type ReplyVoiceMsg struct {
	ReplyMsgBase
	MediaId string
}

func NewReplyVoiceMsg(toUserName, fromUserName, mediaId string) *ReplyVoiceMsg{
	return &ReplyVoiceMsg{
		ReplyMsgBase:newMsgBase(toUserName, fromUserName),
		MediaId:mediaId,
	}
}

func (msg *ReplyVoiceMsg) ToXML() []byte {
	return executeTmpl(RTN_VOICE, msg)
}

// ---------------- reply video mesage -----------------
type ReplyVideoMsg struct {
	ReplyMsgBase
	MediaId string
	Title   string
	Desc    string
}

func NewReplyVideoMsg(toUserName, fromUserName, mediaId, title, desc string) *ReplyVideoMsg {
	return &ReplyVideoMsg{
		ReplyMsgBase:newMsgBase(toUserName, fromUserName),
		MediaId:mediaId,
		Title:title,
		Desc:desc,
	}
}

func (msg *ReplyVideoMsg) ToXML() []byte {
	return executeTmpl(RTN_VIDEO, msg)
}

// ----------------- reply music message -----------------
type ReplyMusicMsg struct {
	ReplyMsgBase
	ThumbMediaId string
	MusicUrl     string
	HQMusicUrl   string
	Title        string
	Desc         string
}

func NewReplyMusicMsg(toUserName, fromUserName, thumbMediaId, musicUrl, HQMusicUrl, title, desc string) *ReplyMusicMsg {
	return &ReplyMusicMsg{
		ReplyMsgBase:newMsgBase(toUserName, fromUserName),
		ThumbMediaId:thumbMediaId,
		MusicUrl:musicUrl,
		HQMusicUrl:HQMusicUrl,
		Title:title,
		Desc:desc,
	}
}

func (msg *ReplyMusicMsg) ToXML() []byte {
	return executeTmpl(RTN_MUSIC, msg)
}

// ------------ reply news message ---------------
type NewsArticle struct {
     Title  string
     Desc   string
     PicUrl string
     Url    string
}

func NewNewsArticle(title, desc, picUrl, url string) *NewsArticle {
	return &NewsArticle {
		Title:title,
		Desc:desc,
		PicUrl:picUrl,
		Url:url,
	}
}

type ReplyNewsMsg struct {
	ReplyMsgBase
	Articles []*NewsArticle
}

func NewReplyNewsMsg(toUserName, fromUserName string, articles []*NewsArticle) *ReplyNewsMsg {
	return &ReplyNewsMsg{
		ReplyMsgBase:newMsgBase(toUserName, fromUserName),
		Articles:articles,
	}
}

func (msg *ReplyNewsMsg) ToXML() []byte {
	return executeTmpl(RTN_NEWS, msg)
}
