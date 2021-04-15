// 微信支付api v3 证书相关API接口
package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/perlyna/wechatpay/model"
)

const certificatesURL = `https://api.mch.weixin.qq.com/v3/certificates`

// GetCertificatesContext 获取平台证书列表
// 文档链接: https://pay.weixin.qq.com/wiki/doc/apiv3_partner/wechatpay/wechatpay5_1.shtml
func GetCertificates(ctx context.Context, hc *http.Client, credential Credential) ([]model.CertificateInfo, error) {
	validator := &WechatPayValidator{&WechatPayDefaultVerifier{}}
	body, err := Get(ctx, hc, credential, validator, certificatesURL)
	if err != nil {
		return nil, err
	}
	reply := model.CertificateReply{}
	if err = json.Unmarshal(body, &reply); err != nil {
		return nil, err
	}
	return reply.Data, nil
}
