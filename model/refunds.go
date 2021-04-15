package model

import "time"

// RefundsAmount 退款金额信息
type RefundsAmount struct {
	Refund   int    `json:"refund"`   // 退款金额, 不能超过原订单支付金额
	Total    int    `json:"total"`    // 原支付交易的订单总金额
	Currency string `json:"currency"` // 目前只支持人民币：CNY

	PayerTotal       int `json:"payer_total,omitempty"`       // 现金支付金额，单位为分
	PayerRefund      int `json:"payer_refund,omitempty"`      // 退款给用户的金额，不包含所有优惠券金额。
	SettlementRefund int `json:"settlement_refund,omitempty"` // 应结退款金额;去掉非充值代金券退款金额后的退款金额，单位为分，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额。
	SettlementTotal  int `json:"settlement_total,omitempty"`  // 应结订单金额;应结订单金额=订单金额-免充值代金券金额，应结订单金额<=订单金额，单位为分。
	DiscountRefund   int `json:"discount_refund,omitempty"`   // 优惠退款金额; 优惠退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金，说明详见代金券或立减优惠，单位为分
}

// RefundsGoodsDetail 退款商品
type RefundsGoodsDetail struct {
	MerchantGoodsID  string `json:"merchant_goods_id"`            // 商户侧商品编码
	WechatpayGoodsID string `json:"wechatpay_goods_id,omitempty"` // 微信侧商品编码
	GoodsName        string `json:"goods_name,omitempty"`         // 商品名称
	UnitPrice        int    `json:"unit_price"`                   // 商品单价,单位为分
	RefundAmount     int    `json:"refund_amount"`                // 商品数量,商品退款金额，单位为分
	RefundQuantity   int    `json:"refund_quantity"`              // 商品退款数量
}

// RefundsReq 退款请求参数
type RefundsReq struct {
	TransactionID string              `json:"transaction_id,omitempty"` // 微信支付订单号
	OutTradeNo    string              `json:"out_trade_no,omitempty"`   // 商户订单号
	OutRefundNo   string              `json:"out_refund_no"`            // 户系统内部的退款单号，商户系统内部唯一
	Reason        string              `json:"reason,omitempty"`         // 退款原因
	NotifyURL     string              `json:"notify_url,omitempty"`     // 退款结果回调url
	FundsAccount  string              `json:"funds_account,omitempty"`  // 退款资金来源, 枚举值：AVAILABLE：可用余额账户
	Amount        RefundsAmount       `json:"amount"`                   // 金额信息
	GoodsDetail   *RefundsGoodsDetail `json:"goods_detail,omitempty"`   // 退款商品
}

type PromotionDetail struct {
	PromotionID  string                `json:"promotion_id"`           // 券ID
	Scope        string                `json:"scope"`                  // 优惠范围
	Type         string                `json:"type"`                   // 优惠类型
	Amount       int                   `json:"amount"`                 // 优惠券面额
	RefundAmount int                   `json:"refund_amount"`          // 优惠退款金额
	GoodsDetails *[]RefundsGoodsDetail `json:"goods_detail,omitempty"` // 商品列表
}

// RefundsOrder 退款订单信息
type RefundsOrder struct {
	RefundID            string          `json:"refund_id"`             // 微信支付退款号
	OutRefundNo         string          `json:"out_refund_no"`         // 商户退款单号
	TransactionID       string          `json:"transaction_id"`        // 微信支付订单号
	OutTradeNo          string          `json:"out_trade_no"`          // 商户订单号
	Channel             string          `json:"channel"`               // 退款渠道
	UserReceivedAccount string          `json:"user_received_account"` // 退款入账账户
	SuccessTime         *time.Time      `json:"success_time"`          // 退款成功时间
	CreateTime          time.Time       `json:"create_time"`           // 退款创建时间
	Status              string          `json:"status"`                // 退款状态
	FundsAccount        string          `json:"funds_account"`         // 资金账户
	Amount              RefundsAmount   `json:"amount"`                // 金额信息
	PromotionDetail     PromotionDetail `json:"promotion_detail"`      // 优惠退款信息
}
