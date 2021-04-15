// 微信支付api v3 签名器
package core

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"
)

// SignatureResult 签名结果
type SignatureResult struct {
	MchCertificateSerialNo string // 商户序列号
	Signature              string // 签名
}

// Signer 签名生成器
type Signer interface {
	GetName() string                                                    // 获取签名器的名称
	GetType() string                                                    // 获取签名器的类型
	GetVersion() string                                                 // 获取签名器的版本
	Sign(ctx context.Context, message string) (*SignatureResult, error) // 对信息进行签名
}

// Sha256WithRSASigner Sha256WithRSA 签名器
type SHA256WithRSASigner struct {
	MchCertificateSerialNo string          // 商户证书序列号
	PrivateKey             *rsa.PrivateKey // 商户私钥
}

// GetName 获取签名器的名称
func (s *SHA256WithRSASigner) GetName() string {
	return "SHA256withRSA"
}

// 获取签名器的类型
func (s *SHA256WithRSASigner) GetType() string {
	return "PRIVATEKEY"
}

// 获取签名器的版本
func (s *SHA256WithRSASigner) GetVersion() string {
	return "1.0"
}

// 对信息使用Sha256WithRsa的方式进行签名
func (s *SHA256WithRSASigner) Sign(ctx context.Context, message string) (*SignatureResult, error) {
	if s.PrivateKey == nil {
		return nil, fmt.Errorf("you must set privatekey to use Sha256WithRSASigner")
	}
	if strings.TrimSpace(s.MchCertificateSerialNo) == "" {
		return nil, fmt.Errorf("you must set mch certificate serial no to use Sha256WithRSASigner")
	}
	h := crypto.Hash.New(crypto.SHA256)
	if _, err := h.Write([]byte(message)); err != nil {
		return nil, err
	}
	hashed := h.Sum(nil)
	signatureByte, err := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}
	ret := &SignatureResult{MchCertificateSerialNo: s.MchCertificateSerialNo}
	ret.Signature = base64.StdEncoding.EncodeToString(signatureByte)
	return ret, nil
}
