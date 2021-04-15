package model

import "time"

type ComplaintReply struct {
	Complaints []Complaint `json:"data"`        // 用户投诉信息详情
	Offset     uint32      `json:"offset"`      // 分页开始位置
	Limit      uint32      `json:"limit"`       // 分页大小
	TotalCount uint64      `json:"total_count"` // 投诉总条数
}

// ComplaintOrderInfo 投诉单关联订单信息
type ComplaintOrderInfo struct {
	TransactionID string `json:"transaction_id"` // 微信订单号
	OutTradeNo    string `json:"out_trade_no"`   // 商户订单号
	Amount        int    `json:"amount"`         // 订单金额
}

// Complaint 投诉单详情
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter10_2_13.shtml
type Complaint struct {
	ComplaintID           string               `json:"complaint_id"`            // 投诉单号
	ComplaintTime         time.Time            `json:"complaint_time"`          // 投诉时间
	ComplaintDetail       string               `json:"complaint_detail"`        // 投诉详情
	ComplaintedMchID      string               `json:"complainted_mchid"`       // 投诉商户号
	ComplaintState        string               `json:"complaint_state"`         // 投诉单状态;PENDING：待处理;PROCESSING：处理中;PROCESSED：已处理完成
	PayerPhone            string               `json:"payer_phone"`             // 投诉人联系方式
	PayerOpenID           string               `json:"payer_openid"`            // 投诉人openid
	Order                 []ComplaintOrderInfo `json:"complaint_order_info"`    // 投诉单关联订单信息
	ComplaintFullRefunded bool                 `json:"complaint_full_refunded"` // 投诉单是否已全额退款
	IncomingUserResponse  bool                 `json:"incoming_user_response"`  // 是否有待回复的用户留言
	UserComplaintTimes    int                  `json:"user_complaint_times"`    // 用户投诉次数。用户首次发起投诉记为1次，用户每有一次继续投诉就加1
}

// ComplaintEvent 投诉通知回调事件
type ComplaintEvent struct {
	ID           string    `json:"id"`            // 通知ID
	CreateTime   time.Time `json:"create_time"`   // 通知创建时间
	EventType    string    `json:"event_type"`    // 通知类型; COMPLAINT.CREATE:产生新投诉; COMPLAINT. STATE_CHANGE:投诉状态变化
	ResourceType string    `json:"resource_type"` // 通知的资源数据类型，支付成功通知为encrypt-resource
	Summary      string    `json:"summary"`       // 回调摘要
	Resource     struct {  // 通知资源数据
		Algorithm      string `json:"algorithm"`       // 加密算法类型,目前只支持AEAD_AES_256_GCM
		Ciphertext     string `json:"ciphertext"`      // Base64编码后的开启/停用结果数据密文
		OriginalType   string `json:"original_type"`   // Base64编码后的开启/停用结果数据密文
		AssociatedData string `json:"associated_data"` // 附加数据
		Nonce          string `json:"nonce"`           // 加密使用的随机串
	} `json:"resource"`

	// 通知资源数据 解密
	ComplaintID string `json:"complaint_id"` // 投诉单号
	ActionType  string `json:"action_type"`  // 触发本次投诉通知回调的具体动作类型
}
