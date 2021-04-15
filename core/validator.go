package core

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Validator 回包校验器
type Validator interface {
	Validate(ctx context.Context, body []byte, header http.Header) error // 对http回包进行校验
}

// WechatPayNullValidator 回包校验器
type WechatPayNullValidator struct {
}

// Validate 使用验证器对回包进行校验
func (validator *WechatPayNullValidator) Validate(ctx context.Context, body []byte, header http.Header) error {
	return nil
}

// WechatPayValidator 回包校验器
type WechatPayValidator struct {
	Verifier Verifier // 验证器
}

// Validate 使用验证器对回包进行校验
func (validator *WechatPayValidator) Validate(ctx context.Context, body []byte, header http.Header) error {
	if validator.Verifier == nil {
		return fmt.Errorf("you must init WechatPayValidator with auth.Verifier")
	}
	err := validateParameters(header)
	if err != nil {
		return err
	}
	message, err := buildMessage(body, header)
	if err != nil {
		return err
	}
	// 微信支付回包平台序列号
	serialNumber := strings.TrimSpace(header.Get(WechatPaySerial))
	// 微信支付回包签名信息
	signature, err := base64.StdEncoding.DecodeString(strings.TrimSpace(header.Get(WechatPaySignature)))
	if err != nil {
		return fmt.Errorf("base64 decode string wechat pay signature err:%s", err.Error())
	}
	err = validator.Verifier.Verify(ctx, serialNumber, message, string(signature))
	if err != nil {
		return fmt.Errorf("validate verify fail serial=%s request-id=%s err=%s", serialNumber,
			strings.TrimSpace(header.Get(RequestID)), err)
	}
	return nil
}

func validateParameters(header http.Header) (err error) {
	// 微信支付回包请求ID
	requestID := strings.TrimSpace(header.Get(RequestID))
	if requestID == "" {
		return fmt.Errorf("empty %s", RequestID)
	}
	// 微信支付回包平台序列号
	if strings.TrimSpace(header.Get(WechatPaySerial)) == "" {
		return fmt.Errorf("empty %s, request-id=[%s]", WechatPaySerial, requestID)
	}
	// 微信支付回包签名信息
	if strings.TrimSpace(header.Get(WechatPaySignature)) == "" {
		return fmt.Errorf("empty %s, request-id=[%s]", WechatPaySignature, requestID)
	}
	// 微信支付回包时间戳
	if strings.TrimSpace(header.Get(WechatPayTimestamp)) == "" {
		return fmt.Errorf("empty %s, request-id=[%s]", WechatPayTimestamp, requestID)
	}
	// 微信支付回包随机字符串
	if strings.TrimSpace(header.Get(WechatPayNonce)) == "" {
		return fmt.Errorf("empty %s, request-id=[%s]", WechatPayNonce, requestID)
	}
	timeStampStr := strings.TrimSpace(header.Get(WechatPayTimestamp))
	timeStamp, err := strconv.Atoi(timeStampStr)
	if err != nil {
		return fmt.Errorf("invalid timestamp:[%s] request-id=[%s] err:[%v]", timeStampStr, requestID, err)
	}
	if math.Abs(float64(timeStamp)-float64(time.Now().Unix())) >= FiveMinute {
		return fmt.Errorf("timestamp=[%d] expires, request-id=[%s]", timeStamp, requestID)
	}
	return nil
}

func buildMessage(body []byte, header http.Header) (message string, err error) {
	timeStamp := strings.TrimSpace(header.Get(WechatPayTimestamp))
	nonce := strings.TrimSpace(header.Get(WechatPayNonce))
	message = fmt.Sprintf("%s\n%s\n%s\n", timeStamp, nonce, body)
	return message, nil
}
