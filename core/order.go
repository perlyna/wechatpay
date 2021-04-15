// 微信支付api v3 订单相关API接口
package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/perlyna/wechatpay/model"
)

// OrderQuery 查询订单API
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_2.shtml
func OrderQuery(ctx context.Context, hc *http.Client, reqURL string, credential Credential, validator Validator) (model.TradeQuery, error) {
	var tradeQuery model.TradeQuery
	body, err := Get(ctx, hc, credential, validator, reqURL)
	if err != nil {
		return tradeQuery, err
	}
	err = json.Unmarshal(body, &tradeQuery)
	return tradeQuery, err
}
