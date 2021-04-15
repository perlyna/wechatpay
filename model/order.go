package model

import "time"

// Amount 订单金额
type Amount struct {
	Total         int    `json:"total"`              // 订单总金额，单位为分。
	PayerTotal    int    `json:"payer_total"`        // 用户支付金额，单位为分。
	Currency      string `json:"currency,omitempty"` // 货币类型; 境内商户号仅支持人民币(CNY)。
	PayerCurrency string `json:"payer_currency"`     // 用户支付币种
}

// Payer 支付者
type Payer struct {
	OpenID string `json:"openid"` // 用户标识
}

// GoodsDetail 商品信息

type GoodsDetail struct {
	MerchantGoodsID  string `json:"merchant_goods_id"`            // 商户侧商品编码
	WechatpayGoodsID string `json:"wechatpay_goods_id,omitempty"` // 微信侧商品编码
	GoodsName        string `json:"goods_name,omitempty"`         // 商品名称
	Quantity         string `json:"quantity"`                     // 商品数量
	UnitPrice        string `json:"unit_price"`                   // 商品单价
}

// Detail 优惠功能
// CostPrice:
// 1、商户侧一张小票订单可能被分多次支付，订单原价用于记录整张小票的交易金额。
// 2、当订单原价与支付金额不相等，则不享受优惠。
// 3、该字段主要用于防止同一张小票分多次支付，以享受多次优惠的情况，正常支付订单不必上传此参数。
type Detail struct {
	CostPrice   int           `json:"cost_price,omitempty"`   // 订单原价
	InvoiceID   string        `json:"invoice_id,omitempty"`   // 商品小票ID
	GoodsDetail []GoodsDetail `json:"goods_detail,omitempty"` // 单品列表信息
}

// StoreInfo 商户门店信息
type StoreInfo struct {
	ID       string `json:"id"`                  // 门店编号
	Name     string `json:"name,omitempty"`      // 门店名称
	AreaCode string `json:"area_code,omitempty"` // 地区编码
	Address  string `json:"address,omitempty"`   // 详细地址
}

// SceneInfo 场景信息
type SceneInfo struct {
	PayerClientIP string    `json:"payer_client_ip"`      // 用户终端IP
	DeviceID      string    `json:"device_id,omitempty"`  // 商户端设备号
	StoreInfo     StoreInfo `json:"store_info,omitempty"` // 商户门店信息
}

// SettleInfo 结算信息
type SettleInfo struct {
	ProfitSharing bool `json:"profit_sharing,omitempty"` // 是否指定分账
}

// UnifiedOrder 统一下单请求参数
type UnifiedOrder struct {
	AppID       string      `json:"appid"`                 // 应用ID
	MchID       string      `json:"mchid"`                 // 商户号
	Description string      `json:"description"`           // 商品描述
	OutTradeNo  string      `json:"out_trade_no"`          // 商户订单号; 商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一
	TimeExpire  string      `json:"time_expire,omitempty"` // 交易结束时间; 订单失效时间，遵循rfc3339标准格式，格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，
	Attach      string      `json:"attach,omitempty"`      // 附加数据
	NotifyURL   string      `json:"notify_url"`            // 通知地址; 通知URL必须为直接可访问的URL，不允许携带查询串。
	GoodsTag    string      `json:"goods_tag,omitempty"`   // 订单优惠标记
	Amount      Amount      `json:"amount"`                // 订单金额
	Payer       *Payer      `json:"payer,omitempty"`       // 支付者信息
	SceneInfo   *SceneInfo  `json:"scene_info,omitempty"`  // 场景信息
	SettleInfo  *SettleInfo `json:"settle_info,omitempty"` // 结算信息
}

// PromotionGoods 优惠的商品信息
type PromotionGoods struct {
	GoodsID        string `json:"goods_id"`        // 商户侧商品编码
	Quantity       int    `json:"quantity"`        // 商品数量
	UnitPrice      int    `json:"unit_price"`      // 商品单价
	DiscountAmount int    `json:"discount_amount"` // 商品优惠金额
	GoodsRemark    string `json:"goods_remark"`    // 商品名称
}

// Promotion 优惠功能, 享受优惠时返回该字段
type Promotion struct {
	CouponID            string           `json:"coupon_id"`            // 券ID
	Name                string           `json:"name"`                 // 优惠名称
	Scope               string           `json:"scope"`                // 优惠范围
	Type                string           `json:"type"`                 // 优惠类型
	Amount              int              `json:"amount"`               // 优惠券面额
	StockID             string           `json:"stock_id"`             // 活动ID
	WechatpayContribute int              `json:"wechatpay_contribute"` // 微信出资
	MerchantContribute  int              `json:"merchant_contribute"`  // 商户出资
	OtherContribute     int              `json:"other_contribute"`     // 其他出资
	Currency            string           `json:"currency"`             // 优惠币种
	PromotionGoods      []PromotionGoods `json:"goods_detail"`         // 商品列表
}

// TradeQuery 交易订单
type TradeQuery struct {
	AppID           string       `json:"appid"`                      // 应用ID
	MchID           string       `json:"mchid"`                      // 商户号
	OutTradeNo      string       `json:"out_trade_no"`               // 商户订单号
	TransactionID   string       `json:"transaction_id"`             // 微信支付订单号
	TradeType       string       `json:"trade_type"`                 // 交易类型
	TradeState      string       `json:"trade_state"`                // 交易状态
	TradeStateDesc  string       `json:"trade_state_desc"`           // 交易状态描述
	BankType        string       `json:"bank_type"`                  // 付款银行
	Attach          string       `json:"attach"`                     // 附加数据
	SuccessTime     time.Time    `json:"success_time"`               // 支付完成时间
	Payer           Payer        `json:"payer"`                      // 支付者信息
	Amount          Amount       `json:"amount"`                     // 订单金额信息，当支付成功时返回该字段
	SceneInfo       *SceneInfo   `json:"scene_info,omitempty"`       // 支付场景描述
	PromotionDetail *[]Promotion `json:"promotion_detail,omitempty"` // 优惠功能, 享受优惠时返回该字段
}
