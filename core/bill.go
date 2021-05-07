package core

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/perlyna/wechatpay/model"
)

// DownloadBill 下载账单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_8.shtml
func DownloadBill(ctx context.Context, hc *http.Client, credential Credential, bill model.Bill) ([]byte, error) {
	body, err := Get(ctx, hc, credential, WithoutValidator, bill.DownloadURL)
	if err != nil {
		return nil, err
	}
	if bill.TarType == "GZIP" {
		if reader, err := gzip.NewReader(bytes.NewBuffer(body)); err == nil {
			if body1, err := ioutil.ReadAll(reader); err == nil {
				body = body1
			}
			reader.Close()
		}
	}
	var h hash.Hash
	switch bill.HashType {
	case "SHA1":
		h = sha1.New()
	default:
		return body, nil
	}
	h.Write(body)
	if hashVal := hex.EncodeToString(h.Sum(nil)); hashVal != bill.HashValue {
		err = fmt.Errorf("校验值不对 [%s] -> [%s]", hashVal, bill.HashValue)
	}
	return body, err
}

const tradebillURL = `https://api.mch.weixin.qq.com/v3/bill/tradebill`

// TradeBill 申请交易账单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_6.shtml
func TradeBill(ctx context.Context, hc *http.Client, credential Credential, validator Validator, date time.Time, billType, tarType string) ([]byte, error) {
	v := url.Values{}
	v.Set("bill_date", date.Format("2006-01-02"))
	if billType == "" {
		billType = "ALL"
	}
	v.Set("bill_type", billType)
	if tarType != "" {
		v.Set("tar_type", tarType)
	}
	body, err := Get(ctx, hc, credential, validator, tradebillURL+"?"+v.Encode())
	if err != nil {
		return body, err
	}
	bill := model.Bill{TarType: tarType}
	if err = json.Unmarshal(body, &bill); err != nil {
		return body, err
	}
	return DownloadBill(ctx, hc, credential, bill)
}

const fundflowillURL = `https://api.mch.weixin.qq.com/v3/bill/fundflowbill`

// FundflowBill 申请资金账单
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_7.shtml
func FundflowBill(ctx context.Context, hc *http.Client, credential Credential, validator Validator, date time.Time, accountType, tarType string) ([]byte, error) {
	v := url.Values{}
	v.Set("bill_date", date.Format("2006-01-02"))
	if accountType != "" {
		v.Set("account_type", accountType)
	}
	if tarType != "" {
		v.Set("tar_type", tarType)
	}
	body, err := Get(ctx, hc, credential, validator, fundflowillURL+"?"+v.Encode())
	if err != nil {
		return body, err
	}
	bill := model.Bill{TarType: tarType}
	if err = json.Unmarshal(body, &bill); err != nil {
		return body, err
	}
	return DownloadBill(ctx, hc, credential, bill)
}
