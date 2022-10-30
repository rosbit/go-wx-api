package ceord

import (
	"encoding/json"
)

// ------------ 订单 -------------
// 订单
type Order struct {
	CreateTime int64 `json:"create_time"` // 秒级时间戳
	UpdateTime int64 `json:"update_time"` // 秒级时间戳
	OrderId string `json:"order_id"` // 订单号
	Status int16 `json:"status"` // 订单状态，10: 待付款, 20: 待发货, 21: 部分发货, 30: 待收货, 100: 完成, 190: 商品超卖商家取消订单, 200: 全部商品售后之后，订单取消, 250	未付款用户主动取消或超时未付款订单自动取消
	OpenId string `json:"openid"` // 买家身份标识
	UnionId string `json:"unionid"`
	OrderDetail `json:"order_detail"` // 订单详细数据信息
	AfterSaleDetail `json:"aftersale_detail"` // 售后信息
}

// 订单详细数据
type OrderDetail struct {
	ProductInfos []ProductInfo `json:"product_infos"` // 商品列表
	PriceInfo `json:"price_info"` // 价格信息
	PayInfo `json:"pay_info"` // 支付信息
	DeliveryInfo `json:"delivery_info"` // 配送信息
	CouponInfo `json:"coupon_info"` // 优惠券信息
	ExtInfo `json:"ext_info"` // 额外信息
}

// 商品信息
type ProductInfo struct {
	ProductId uint64 `json:"product_id"` // 商品spuid
	SkuId uint64 `json:"sku_id"` // 商品skuid
	ThumbImg string `json:"thumb_img"` // sku小图
	SkuCount int `json:"sku_cnt"` // sku数量
	SalePrice int32 `json:"sale_price"` // 售卖价格（单位：分）
	Title string `json:"title"` // 商品标题
	OnAfterSaleSkuCount int `json:"on_aftersale_sku_cnt"` // 正在售后/退款流程中的 sku 数量
	FinishAfterSaleSkuCount int `json:"finish_aftersale_sku_cnt"` // 完成售后/退款的 sku 数量
	SkuCode string `json:"sku_code"` // 商品编码
	MarketPrice int32 `json:"market_price"` // 市场价格（单位：分）
	SkuAttrs []AttrInfo `json:"sku_attrs"` // sku属性
	RealPrice int32 `json:"real_price"` // sku实付价格
	OutProductId string `json:"out_product_id"` // 商品外部spuid
	OutSkuId string `json:"out_sku_id"` // 商品外部skuid
}

// sku属性
type AttrInfo struct {
	AttrKey   string `json:"attr_key"` // 属性键（属性自定义用）
	AttrValue string `json:"attr_value"` // 属性值（属性自定义用）
}

// 支付信息
type PayInfo struct {
	PrepayId string `json:"prepay_id"` // 预支付id
	PrepayTime int64 `json:"prepay_time"` // 预支付时间，秒级时间戳
	PayTime int64 `json:"pay_time"` // 支付时间，秒级时间戳
	TransactionId string `json:"transaction_id"` // 支付订单号
}

// 价格信息
type PriceInfo struct {
	ProductPrice int32 `json:"product_price"` // 商品总价，单位为分
	OrderPrice int32 `json:"order_price"` //订单金额，单位为分
	Freight int32 `json:"freight"` // 运费，单位为分
	DiscountedPrice int32 `json:"discounted_price"` // 优惠金额，单位为分
	IsDicounted bool `json:"is_discounted"` // 是否有优惠
	OriginalOrderPrice int32 `json:"original_order_price"` // 订单原始价格，单位为分
	EstimateProductPrice int32 `json:"estimate_product_price"` // 商品预估价格，单位为分
	ChangeDownPrice int32 `json:"change_down_price"` // 改价后降低金额，单位为分
	ChangeFreight int32 `json:"change_freight"` // 改价后运费，单位为分
	IsChangeFreight bool `json:"is_change_freight"` // 是否修改运费
}

// 配送信息
type DeliveryInfo struct {
	AddressInfo `json:"address_info"` // 地址信息
	DeliveryProductInfos []DeliveryProductInfo `json:"delivery_product_info"` // 发货物流信息
	ShipDoneTime int64 `json:"ship_done_time"` // 发货完成时间，秒级时间戳
	DeliverMethod int16 `json:"deliver_method"` // 订单发货方式，0：普通物流；1：虚拟发货，由商品的同名字段决定
}

// 发货物流信息
type DeliveryProductInfo struct {
	WaybillId string `json:"waybill_id"` // 快递单号
	DeliveryId string `json:"delivery_id"` // 快递公司编码
	ProductInfos []FreightProductInfo `json:"product_infos"` //	包裹中商品信息
	DeliveryName string `json:"delivery_name"` // 快递公司名称
	DeliveryTime int64 `json:"delivery_time"` // 发货时间，秒级时间戳
	DeliveryType int16 `json:"deliver_type"` // 配送方式，1:自寄快递, 2:在线签约快递单, 4:在线快递散单
	DeliveryAddress AddressInfo `json:"delivery_address"` // 发货地址
}

// 包裹中商品信息
type FreightProductInfo struct {
	ProductId string `json:"product_id"` // 商品id
	SkuId string `json:"sku_id"` // sku_id
	ProductCount int `json:"product_cnt"` // 商品数量
}

// 发货地址
type AddressInfo struct {
	UserName string `json:"user_name"` // 收货人姓名
	PostalCode string `json:"postal_code"` // 邮编
	ProvinceName string `json:"province_name"` // 省份
	CityName string `json:"city_name"` // 城市
	CountryName string `json:"county_name"` // 区
	DetailInfo string `json:"detail_info"` // 详细地址
	NationalCode string `json:"national_code"` // 国家码
	TelNumber string `json:"tel_number"` // 联系方式
	HouseNumber string `json:"house_number"` // 门牌号码
	VirtualOrderTelNumber string `json:"virtual_order_tel_number"` // 虚拟发货订单联系方式(deliver_method=1时返回)
}

// 优惠券信息
type CouponInfo struct {
	UserCouponId string `json:"user_coupon_id"` // 用户优惠券id
}

// 售后信息
type AfterSaleDetail struct {
	OnAfterSaleOrderCount int `json:"on_aftersale_order_cnt"` // 正在售后流程的售后单数
	AfterSaleOrderList []AfterSaleOrderInfo `json:"aftersale_order_list"` // 售后单列表
}

// 售后单信息
type AfterSaleOrderInfo struct {
	AfterSaleOrderId uint64 `json:"aftersale_order_id"` // 售后单ID
	Status int16 `json:"status"` // 售后单状态(已废弃，请勿使用，售后信息请调用售后接口）
}

// 额外信息
type ExtInfo struct {
	CustomerNotes string `json:"customer_notes"` // 用户备注
	MerchantNotes string `json:"merchant_notes"` // 商家备注
}

// ------------ 售后单 -------------
type AfterSaleOrder struct {
	AfterSaleOrderId string `json:"after_sale_order_id"` // 售后单号
	Status string `json:"status"` // 售后单当前状态
									// USER_CANCELD	用户取消申请
									// MERCHANT_PROCESSING	商家受理中
									// MERCHANT_REJECT_REFUND	商家拒绝退款
									// MERCHANT_REJECT_RETURN	商家拒绝退货退款
									// USER_WAIT_RETURN	待买家退货
									// RETURN_CLOSED	退货退款关闭
									// MERCHANT_WAIT_RECEIPT	待买家退货
									// MERCHANT_OVERDUE_REFUND	商家逾期未退款
									// MERCHANT_REFUND_SUCCESS	退款完成
									// MERCHANT_RETURN_SUCCESS	退货退款完成
									// PLATFORM_REFUNDING	平台退款中
									// PLATFORM_REFUND_FAIL	平台退款失败
									// USER_WAIT_CONFIRM	待用户确认
									// MERCHANT_REFUND_RETRY_FAIL	商家打款失败，客服关闭售后
									// MERCHANT_FAIL	售后关闭
	OpenId string `json:"openid"` // 买家身份标识
	UninoId string `json:"unionid"` // 买家在开放平台的唯一标识符，若当前视频号小店已绑定到微信开放平台帐号下会返回，详见UnionID 机制说明
	ProductInfo AfterSaleProductInfo `json:"product_info"` // 售后相关商品信息
	Details AfterSaleDetails `json:"details"` // 退款详情
	RefundInfo `json:"refund_info"` // 退款信息
	ReturnInfo `json:"return_info"` // 用户退货信息
	MerchantUploadInfo `json:"merchant_upload_info"` // 商家上传的信息
	CreateTime int64 `json:"create_time"` // 售后单创建时间戳
	UpdateTime int64 `json:"update_time"` // 售后单更新时间戳
	Reason string `json:"reason"` // 退款原因
									// INCORRECT_SELECTION	拍错/多拍
									// NO_LONGER_WANT	不想要了
									// NO_EXPRESS_INFO	无快递信息
									// EMPTY_PACKAGE	包裹为空
									// REJECT_RECEIVE_PACKAGE	已拒签包裹
									// NOT_DELIVERED_TOO_LONG	快递长时间未送达
									// NOT_MATCH_PRODUCT_DESC	与商品描述不符
									// QUALITY_ISSUE	质量问题
									// SEND_WRONG_GOODS	卖家发错货
									// THREE_NO_PRODUCT	三无产品
									// FAKE_PRODUCT	假冒产品
									// OTHERS	其它
	RefundResp json.RawMessage `json:"refund_resp"` //	退款结果
	Type string `json:"type"` // 售后类型。REFUND:退款；RETURN:退货退款。
}

// 售后相关商品信息
type AfterSaleProductInfo struct {
	ProductId string `json:"product_id"` // 商品spuid
	SkuId string `json:"sku_id"` // 商品skuid
	Count int `json:"count"` // 售后数量
}

// 退款信息
type RefundInfo struct {
	Amount int32 `json:"amount"` // 退款金额（分）
}

// 用户退货信息
type ReturnInfo = json.RawMessage

// 商家上传的信息
type MerchantUploadInfo struct {
	RejectReason string `json:"reject_reason"` // 拒绝原因
	RefundCertificates []json.RawMessage `json:"refund_certificates"` // 退款凭证
}

// 退款详情
type AfterSaleDetails struct {
	Desc string `json:"desc"` // 售后描述
	ReceiveProduct bool `json:"receive_product"` //	是否已经收到货
	CancelTime int64 `json:"cancel_time"` // 取消售后时间
	ProveImgs []string `json:"prove_imgs"` // 举证图片列表
	TelNumber string `json:"tel_number"` // 联系电话
}

