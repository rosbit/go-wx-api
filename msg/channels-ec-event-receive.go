/**
 * 视频号小店事件的结构定义和解析
 * Rosbit Xu
 */
package wxmsg

const (
	// MT_EVENT = "event" // special message type

	// event type
	ET_ORDER_CANCEL = "channels_ec_order_cancel"
	ET_ORDER_PAY    = "channels_ec_order_pay"
	ET_ORER_CONFIRM = "channels_ec_order_confirm"
	ET_ORDER_SETTLE = "channels_ec_order_settle"
	ET_AFTERSAL_UPDATE = "channels_ec_aftersale_update"
)

type ReceivedJSONEvent interface{}

// ------------- event ---------------
type ChannelsEcEvent struct {
	ToUserName string
	FromUserName string
	CreateTime int
	MsgType string
	Event string
}

// ---- event channels_ec_order_cancel ----
type OrderCancelEvent struct {
	ChannelsEcEvent
	OrderInfo struct {
		OrderId uint64 `json:"order_id"`
		CancelType int16 `json:"cancel_type"`
	} `json:"order_info"`
}

// ---- event channels_ec_order_pay ----
type OrderPayEvent struct {
	ChannelsEcEvent
	OrderInfo struct {
		OrderId uint64 `json:"order_id"`
		PayTime int64 `json:"pay_time"`
	} `json:"order_info"`
}

// ---- event channels_ec_order_confirm ----
type OrderConfirmEvent struct {
	ChannelsEcEvent
	OrderInfo struct {
		OrderId uint64 `json:"order_id"`
		ConfirmType int16 `json:"confirm_type"`
	} `json:"order_info"`
}

// ---- event channels_ec_order_settle ----
type OrderSettleEvent struct {
	ChannelsEcEvent
	OrderInfo struct {
		OrderId uint64 `json:"order_id"`
		SettleTime int64 `json:"settle_time"`
	} `json:"order_info"`
}

// ---- event channels_ec_aftersale_update ----
type AftersaleUpdateEvent struct {
	ChannelsEcEvent
	AftersaleOrderInfo struct {
		Status string `json:"status"`
		AfterSaleOrderid string `json:"after_sale_order_id"`
	} `json:"finder_shop_aftersale_status_update"`
}
