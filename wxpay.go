package wechatpay

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/perlyna/wechatpay/core"
	"github.com/perlyna/wechatpay/model"
	"github.com/perlyna/wechatpay/util"
)

// WechatPay 微信支付SDK
type WechatPay struct {
	mchID                   string                       // 微信商户号
	apiv3Secret             string                       // 商户号 API Secret
	privateKey              *rsa.PrivateKey              // 商户私钥 apiclient_key.pem
	certificates            map[string]*x509.Certificate // 商户密钥 apiclient_cert.pem
	certificateSerialNumber string                       // 商户密钥证书序列号
	credential              core.Credential              // 授权信息生成器
	validator               core.Validator               // 签名校验相关接口

	NotifyURL string       // 支付通知地址
	Client    *http.Client // http client
}

// New 创建微信支付模块
func New(mchid string, apiv3Secret string, privateKey *rsa.PrivateKey, certificate *x509.Certificate) *WechatPay {
	serialNumber := strings.ToUpper(hex.EncodeToString(certificate.SerialNumber.Bytes()))
	signer := &core.SHA256WithRSASigner{MchCertificateSerialNo: serialNumber, PrivateKey: privateKey}
	certificates := make(map[string]*x509.Certificate)
	certificates[serialNumber] = certificate
	verifier := &core.WechatPayVerifier{Certificates: certificates}
	config := &WechatPay{
		mchID:                   mchid,
		apiv3Secret:             apiv3Secret,
		privateKey:              privateKey,
		certificates:            certificates,
		certificateSerialNumber: serialNumber,
		credential:              &core.WechatPayCredentials{Signer: signer, MchID: mchid},
		validator:               &core.WechatPayValidator{Verifier: verifier},
		Client:                  http.DefaultClient,
	}
	return config
}

// UpdateCertificates 更新商户当前可用的平台证书列表
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/wechatpay/wechatpay5_1.shtml
func (p *WechatPay) UpdateCertificates() error {
	certs, err := core.GetCertificates(context.Background(), p.Client, p.credential)
	if err != nil {
		return err
	}
	for _, cert := range certs {
		if cert.ExpireTime.Before(time.Now()) { // 证书已过期
			continue
		}
		if _, ok := p.certificates[cert.SerialNo]; ok { // 证书已存在
			continue
		}
		rawCert, err := util.DecryptToByte(p.apiv3Secret, cert.EncryptCertificate.AssociatedData,
			cert.EncryptCertificate.Nonce, cert.EncryptCertificate.Ciphertext)
		if err != nil {
			return err
		}
		certificate, err := util.LoadCertificate(rawCert)
		if err != nil {
			return err
		}
		serialNumber := strings.ToUpper(hex.EncodeToString(certificate.SerialNumber.Bytes()))
		p.certificates[serialNumber] = certificate
	}
	return nil
}

// OrderQueryByTransactions 微信支付订单号查询
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_2.shtml
func (p *WechatPay) OrderQueryByTransactions(ctx context.Context, transactionID string) (model.TradeQuery, error) {
	reqURL := "https://api.mch.weixin.qq.com/v3/pay/transactions/id/" + transactionID + "?mchid=" + p.mchID
	return core.OrderQuery(ctx, p.Client, reqURL, p.credential, p.validator)
}

// OrderQueryByOutTradeNo 商户订单号查询
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_2.shtml
func (p *WechatPay) OrderQueryByOutTradeNo(ctx context.Context, outTradeNo string) (model.TradeQuery, error) {
	reqURL := "https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/" + outTradeNo + "?mchid=" + p.mchID
	return core.OrderQuery(ctx, p.Client, reqURL, p.credential, p.validator)
}

// RefundByTransactions 微信支付订单号申请退款
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_9.shtml
func (p *WechatPay) RefundByTransactions(ctx context.Context, transactionID string, amount int) (model.RefundsOrder, error) {
	var refundsReq model.RefundsReq
	tradeQuery, err := p.OrderQueryByTransactions(ctx, transactionID)
	if err != nil {
		return model.RefundsOrder{}, err
	}
	if amount == 0 {
		amount = tradeQuery.Amount.PayerTotal
	}
	refundsReq.TransactionID = transactionID
	refundsReq.OutRefundNo = tradeQuery.OutTradeNo + strconv.FormatInt(time.Now().Unix(), 36)
	refundsReq.Amount.Currency = tradeQuery.Amount.PayerCurrency
	refundsReq.Amount.Total = tradeQuery.Amount.Total
	refundsReq.Amount.Refund = amount
	return core.Refunds(ctx, p.Client, refundsReq, p.credential, p.validator)
}

// RefundByOutTradeNo 商户订单号申请退款
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_9.shtml
func (p *WechatPay) RefundByOutTradeNo(ctx context.Context, outTradeNo string, amount int) (model.RefundsOrder, error) {
	var refundsReq model.RefundsReq
	tradeQuery, err := p.OrderQueryByOutTradeNo(ctx, outTradeNo)
	if err != nil {
		return model.RefundsOrder{}, err
	}
	if amount == 0 {
		amount = tradeQuery.Amount.PayerTotal
	}

	refundsReq.OutTradeNo = outTradeNo
	// refundsReq.TransactionID = tradeQuery.TransactionID
	refundsReq.OutRefundNo = tradeQuery.OutTradeNo + strconv.FormatInt(time.Now().UnixNano(), 36)
	refundsReq.Amount.Currency = tradeQuery.Amount.PayerCurrency
	refundsReq.Amount.Total = tradeQuery.Amount.Total
	refundsReq.Amount.Refund = amount
	return core.Refunds(ctx, p.Client, refundsReq, p.credential, p.validator)
}
