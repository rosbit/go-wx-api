/**
 * 视频号小店事件的缺省处理函数
 */
package wxmsg

// 视频号小店事件处理接口定义
type ChannelsEcEventHandler interface {
	HandleOrderCancelEvent(orderCancelEvent *OrderCancelEvent) []byte
	HandleOrderPayEvent(orderPayEvent *OrderPayEvent) []byte
	HandleOrderConfirmEvent(orderConfirmEvent *OrderConfirmEvent) []byte
	HandleOrderSettleEvent(orderSettleEvent *OrderSettleEvent) []byte
	HandleAftersaleUpdateEvent(aftersaleUpdateEvent *AftersaleUpdateEvent) []byte
}

// 缺省的视频号小店事件处理器，可以通过RegisterChannelsEcEventhandler覆盖
var CEEventHandler ChannelsEcEventHandler = &ChannelsEcEventHandlerAdapter{}

// 缺省的视频号小店事件处理器实现。如果要重新实现某些处理方法，需要
// 1. 在新的结构体中嵌入(embed)该结构体
// 2. 覆盖实现某些方法
// 3. 调用RegisterChannelsEcEventhandler覆盖缺省实现
type ChannelsEcEventHandlerAdapter struct {
}

func (h *ChannelsEcEventHandlerAdapter) HandleOrderPayEvent(orderPayEvent *OrderPayEvent) []byte {
	return SUCCESS_TEXT
}

func (h *ChannelsEcEventHandlerAdapter) HandleOrderCancelEvent(orderCancelEvent *OrderCancelEvent) []byte {
	return SUCCESS_TEXT
}

func (h *ChannelsEcEventHandlerAdapter) HandleOrderConfirmEvent(orderConfirmEvent *OrderConfirmEvent) []byte {
	return SUCCESS_TEXT
}

func (h *ChannelsEcEventHandlerAdapter) HandleOrderSettleEvent(orderSettleEvent *OrderSettleEvent) []byte {
	return SUCCESS_TEXT
}

func (h *ChannelsEcEventHandlerAdapter) HandleAftersaleUpdateEvent(aftersaleUpdateEvent *AftersaleUpdateEvent) []byte {
	return SUCCESS_TEXT
}

type FnJSONEventHandler func(channelsEcEvent ReceivedJSONEvent) []byte

func (p *ChannelsEcEventParser) handleReceivedEvent(channelsEcEvent ReceivedJSONEvent, eventType string) []byte {
	if fn, ok := p.eventHandlers[eventType]; ok {
		return fn(channelsEcEvent)
	}
	return nil
}

func (p *ChannelsEcEventParser) RegisterChannelsEcEventHandler(eventHandler ChannelsEcEventHandler) {
	if eventHandler == nil {
		return
	}

	p.eventHandlers = map[string]FnJSONEventHandler {
		ET_ORDER_CANCEL:  func(channelsEcEvent ReceivedJSONEvent) []byte { return eventHandler.HandleOrderCancelEvent(channelsEcEvent.(*OrderCancelEvent)) },
		ET_ORDER_PAY: func(channelsEcEvent ReceivedJSONEvent) []byte { return eventHandler.HandleOrderPayEvent(channelsEcEvent.(*OrderPayEvent)) },
		ET_ORER_CONFIRM:   func(channelsEcEvent ReceivedJSONEvent) []byte { return eventHandler.HandleOrderConfirmEvent(channelsEcEvent.(*OrderConfirmEvent)) },
		ET_ORDER_SETTLE: func(channelsEcEvent ReceivedJSONEvent) []byte { return eventHandler.HandleOrderSettleEvent(channelsEcEvent.(*OrderSettleEvent)) },
		ET_AFTERSAL_UPDATE: func(channelsEcEvent ReceivedJSONEvent) []byte { return eventHandler.HandleAftersaleUpdateEvent(channelsEcEvent.(*AftersaleUpdateEvent)) },
	}
}

