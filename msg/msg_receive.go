/**
 * 服务号消息/事件的结构定义和解析
 * Rosbit Xu
 */
package wxmsg

import (
	"github.com/beevik/etree"
	"fmt"
	"strconv"
)

const (
	// msgType
	MT_TEXT  = "text"
	MT_IMAGE = "image"
	MT_VOICE = "voice"
	MT_VIDEO = "video"
	MT_SHORTVIDEO = "shortvideo"
	MT_LOCATION   = "location"
	MT_LINK  = "link"
	MT_EVENT = "event" // special message type

	// event type
	ET_VIEW  = "VIEW"
	ET_CLICK = "CLICK"
	ET_SUBSCRIBE   = "subscribe"
	ET_UNSUBSCRIBE = "unsubscribe"
	ET_SCAN        = "SCAN"  // SCAN事件内容和subscribe是一样的，注意和scancode_xxx事件的区分
	ET_WHERE       = "location" // 由location_select菜单引起的事件
	ET_LOCATION    = "LOCATION"
	ET_PIC_SYSPHOTO       = "pic_sysphoto"
	ET_PIC_PHOTO_OR_ALBUM = "pic_photo_or_album"
	ET_PIC_WEIXIN         = "pic_weixin"
	ET_SCANCODE_WAITMSG   = "scancode_waitmsg"
	ET_SCANCODE_PUSH      = "scancode_push"
	ET_MASSSENDJOBFINISH     = "MASSSENDJOBFINISH"
	ET_TEMPLATESENDJOBFINISH = "TEMPLATESENDJOBFINISH"
)

func _getText(el *etree.Element, tagName string) (string, error) {
	t := el.SelectElement(tagName)
	if t == nil {
		return "", fmt.Errorf("%s not found", tagName)
	}
	return t.Text(), nil
}

// 所有接收消息都实现的接口。消息处理可以统一为 ReceivedMsg -> [处理] -> ReplyMsg
type ReceivedMsg interface {
	parse(root *etree.Element)
}

// --------------------- message -------------------
type BaseMsg struct {
	ToUserName string
	FromUserName string
	CreateTime int
	MsgType string
}

func (m *BaseMsg) parse(root *etree.Element) {
	m.ToUserName, _   = _getText(root, "ToUserName")
	m.FromUserName, _ = _getText(root, "FromUserName")
	m.MsgType, _      = _getText(root, "MsgType")
	ct, _            := _getText(root, "CreateTime")
	m.CreateTime, _ = strconv.Atoi(ct)
}

type Msg struct {
	BaseMsg
	MsgId string
}

func (m *Msg) parse(root *etree.Element) {
	m.BaseMsg.parse(root)
	m.MsgId, _ = _getText(root, "MsgId")
}

// ---- text message ----
type TextMsg struct {
	Msg
	Content string
}

func (m *TextMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.Content, _ = _getText(root, "Content")
}

// ---- image message ----
type ImageMsg struct {
	Msg
	PicUrl  string
	MediaId string
}

func (m *ImageMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.PicUrl, _  = _getText(root, "PicUrl")
	m.MediaId, _ = _getText(root, "MediaId")
}

// ---- voice message ----
type VoiceMsg struct {
	Msg
	MediaId string
	Format  string
	Recognition string
}

func (m *VoiceMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.MediaId, _      = _getText(root, "MediaId")
	m.Format, _       = _getText(root, "Format")
	m.Recognition, _  = _getText(root, "Recognition")
}

// ---- video message ----
type VideoMsg struct {
	Msg
	MediaId      string
	ThumbMediaId string
}

func (m *VideoMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.MediaId, _      = _getText(root, "MediaId")
	m.ThumbMediaId, _ = _getText(root, "ThumbMediaId")
}

// ---- location message -----
type LocationMsg struct {
	Msg
	Location_X string
	Location_Y string
	Scale int
	Label string
}

func (m *LocationMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.Location_X, _ = _getText(root, "Location_X")
	m.Location_Y, _ = _getText(root, "Location_Y")
	scale, _   := _getText(root, "Scale")
	m.Scale, _  = strconv.Atoi(scale)
	m.Label, _  = _getText(root, "Label")
}

// ---- link message -----
type LinkMsg struct {
	Msg
	Title string
	Description string
	Url   string
}

func (m *LinkMsg) parse(root *etree.Element) {
	m.Msg.parse(root)
	m.Title, _       = _getText(root, "Title")
	m.Description, _ = _getText(root, "Description")
	m.Url, _         = _getText(root, "Url")
}

// ------------- event ---------------
type EventMsg struct {
	BaseMsg
	Event string
}

func (m *EventMsg) parse(root *etree.Element) {
	m.BaseMsg.parse(root)
	m.Event, _ = _getText(root, "Event")
}

// ---- event CLICK ----
type ClickEvent struct {
	EventMsg
	EventKey string
}

func (m *ClickEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.EventKey, _ = _getText(root, "EventKey")
}

// ---- event VIEW ----
type ViewEvent struct {
	EventMsg
	EventKey string
	MenuId string
}

func (m *ViewEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.EventKey, _ = _getText(root, "EventKey")
	m.MenuId, _   = _getText(root, "MenuId")
}

// ---- event subscribe, unsubscribe, SCAN ----
type SubscribeEvent struct {
	EventMsg
	EventKey string
	Ticket string
}

func (m *SubscribeEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.EventKey, _ = _getText(root, "EventKey")
	m.Ticket, _   = _getText(root, "Ticket")
}

// ---- event pic_sysphoto, pic_photo_or_album, pic_weixin ---
type PhotoEvent struct {
	EventMsg
	EventKey string
	Count int         // SendPicsInfo/Count
	PicList []string  // SendPicsInfo/item/PicMd5Sum
}

func (m *PhotoEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.EventKey, _  = _getText(root, "EventKey")
	sendPicsInfo := root.SelectElement("SendPicsInfo")
	if sendPicsInfo == nil {
		return
	}
	c, _ := _getText(sendPicsInfo, "Count")
	m.Count, _  = strconv.Atoi(c)
	if m.Count <= 0 {
		return
	}
	picList := sendPicsInfo.SelectElement("PicList")
	if picList == nil {
		return
	}
	items := picList.SelectElements("item")
	if items == nil {
		return
	}

	m.PicList = make([]string, len(items))
	for i, item := range items {
		m.PicList[i], _ = _getText(item, "PicMd5Sum")
	}
}

// ---- event location ----- 由location_select菜单引起的事件
type WhereEvent struct {
	EventMsg
	Location_X string
	Location_Y string
	Scale      int
	Label      string
	MsgId      string
}

func (m *WhereEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.Location_X, _ = _getText(root, "Location_X")
	m.Location_Y, _ = _getText(root, "Location_Y")
	scale, _       := _getText(root, "Scale")
	m.Scale, _      = strconv.Atoi(scale)
	m.Label, _      = _getText(root, "Label")
	m.MsgId, _      = _getText(root, "MsgId")
}

// ---- event LOCATION -----
type LocationEvent struct {
	EventMsg
	Latitude  string
	Longitude string
	Precision string
}

func (m *LocationEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.Latitude, _  = _getText(root, "Latitude")
	m.Longitude, _ = _getText(root, "Longitude")
	m.Precision, _ = _getText(root, "Precision")
}

// ---- event scancode_waitmsg, scancode_push ----
type ScancodeEvent struct {
	EventMsg
	EventKey   string
	ScanType   string
	ScanResult string
}

func (m *ScancodeEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.EventKey, _   = _getText(root, "EventKey")
	scanCodeInfo   := root.SelectElement("ScanCodeInfo")
	if scanCodeInfo == nil {
		return
	}
	m.ScanType, _   = _getText(scanCodeInfo, "ScanType")
	m.ScanResult, _ = _getText(scanCodeInfo, "ScanResult")
}

// ---- event MASSSENDJOBFINISH ---- 群发结果推送事件
type ResultListItem struct {
	ArticleIdx            string
	UserDeclareState      string
	AuditState            string
	OriginalArticleUrl    string
	OriginalArticleType   string
	CanReprint            string
	NeedReplaceContent    string
	NeedShowReprintSource string
}

func (m *ResultListItem) parse(root *etree.Element) {
	m.ArticleIdx, _            = _getText(root, "ArticleIdx")
	m.UserDeclareState, _      = _getText(root, "UserDeclareState")
	m.AuditState, _            = _getText(root, "AuditState")
	m.OriginalArticleUrl, _    = _getText(root, "OriginalArticleUrl")
	m.OriginalArticleType, _   = _getText(root, "OriginalArticleType")
	m.CanReprint, _            = _getText(root, "CanReprint")
	m.NeedReplaceContent, _    = _getText(root, "NeedReplaceContent")
	m.NeedShowReprintSource, _ = _getText(root, "NeedShowReprintSource")
}

type MassSentEvent struct {
	EventMsg
	MsgID  string
	Status string
	TotalCount  int
	FilterCount int
	SentCount   int
	ErrorCount  int
	ResultCount int
	ResultList  []ResultListItem
}

func (m *MassSentEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.MsgID, _  = _getText(root, "MsgID")
	m.Status, _ = _getText(root, "Status")

	c, _ := _getText(root, "TotalCount")
	m.TotalCount, _  = strconv.Atoi(c)
	c, _ = _getText(root, "FilterCount")
	m.FilterCount, _ = strconv.Atoi(c)
	c, _ = _getText(root, "SentCount")
	m.SentCount, _   = strconv.Atoi(c)
	c, _ = _getText(root, "ErrorCount")
	m.ErrorCount, _  = strconv.Atoi(c)

	copyrightCheckResult := root.SelectElement("CopyrightCheckResult")
	if copyrightCheckResult == nil {
		return
	}
	c, _ = _getText(copyrightCheckResult, "Count")
	m.ResultCount, _ = strconv.Atoi(c)
	if m.ResultCount <= 0 {
		return
	}

	m.ResultList = make([]ResultListItem, m.ResultCount)
	resultList := copyrightCheckResult.SelectElement("ResultList")
	if resultList == nil {
		return
	}
	items := resultList.SelectElements("item")
	if items == nil {
		return
	}
	for i, item := range items {
		m.ResultList[i].parse(item)
	}
}

// ---- event TEMPLATESENDJOBFINISH ----- 模板消息发送结果推送事件
type TemplateSentEvent struct {
	EventMsg
	MsgId  string
	Status string
}

func (m *TemplateSentEvent) parse(root *etree.Element) {
	m.EventMsg.parse(root)
	m.MsgId, _  = _getText(root, "MsgID")
	m.Status, _ = _getText(root, "Status")
}

