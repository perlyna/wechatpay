// 微信支付api v3 authorization生成器
package core

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Credential Authorization信息生成器
type Credential interface {
	GenerateAuthorizationHeader(ctx context.Context, method, canonicalURL,
		signBody string) (authorization string, err error)
}

// WechatPayCredentials authorization生成器
type WechatPayCredentials struct {
	Signer Signer // 签名器
	MchID  string // 商户号
}

// GenerateAuthorizationHeader  生成http request header 中的authorization信息
func (c *WechatPayCredentials) GenerateAuthorizationHeader(ctx context.Context,
	method, canonicalURL, signBody string) (authorization string, err error) {
	if c.Signer == nil {
		return "", fmt.Errorf("you must init WechatPayCredentials with signer")
	}
	nonce := GenerateNonceStr(32)

	timestamp := time.Now().Unix()
	message := fmt.Sprintf(FormatMessage, method, canonicalURL, timestamp, nonce, signBody)
	signatureResult, err := c.Signer.Sign(ctx, message)
	if err != nil {
		return "", err
	}
	authorization = fmt.Sprintf(HeaderAuthorization, c.MchID, nonce, timestamp,
		signatureResult.MchCertificateSerialNo, signatureResult.Signature)
	return authorization, nil
}

// nonceStr 随机字符串
const nonceStr = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// nonceStrLength 随机字符串长度
const nonceStrLength = len(nonceStr)

// GenerateNonceStr 获取随机字符串
func GenerateNonceStr(length int) string {
	bytes := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		bytes[i] = nonceStr[r.Intn(nonceStrLength)]
	}
	return string(bytes)
}
