package wechatpay

import (
	"context"
	"time"

	"github.com/perlyna/wechatpay/core"
)

// TradeBill 申请交易账单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_6.shtml
func (p *WechatPay) TradeBill(ctx context.Context, date time.Time, billType string) ([]byte, error) {
	return core.TradeBill(ctx, p.Client, p.credential, p.validator, date, billType, "GZIP")
}

// FundflowBill 申请资金账单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_7.shtml
func (p *WechatPay) FundflowBill(ctx context.Context, date time.Time, accountType string) ([]byte, error) {
	return core.FundflowBill(ctx, p.Client, p.credential, p.validator, date, accountType, "GZIP")
}
