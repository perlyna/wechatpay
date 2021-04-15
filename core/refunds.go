package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/perlyna/wechatpay/model"
)

const refundsURL = `https://api.mch.weixin.qq.com/v3/refund/domestic/refunds`

func Refunds(ctx context.Context, hc *http.Client, refundsReq model.RefundsReq, credential Credential, validator Validator) (model.RefundsOrder, error) {
	var refundsOrder model.RefundsOrder
	body, err := Post(ctx, hc, credential, validator, refundsURL, refundsReq)
	if err != nil {
		return refundsOrder, err
	}
	err = json.Unmarshal(body, &refundsOrder)
	return refundsOrder, err
}
