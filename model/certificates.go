package model

import "time"

type CertificateInfo struct {
	EffectiveTime      time.Time `json:"effective_time"` // 证书启用时间
	ExpireTime         time.Time `json:"expire_time"`    // 证书过期时间
	SerialNo           string    `json:"serial_no"`      // 证书序列号
	EncryptCertificate struct {
		Algorithm      string `json:"algorithm"`
		AssociatedData string `json:"associated_data"`
		Ciphertext     string `json:"ciphertext"`
		Nonce          string `json:"nonce"`
	} `json:"encrypt_certificate"` // 证书加密信息
}

type CertificateReply struct {
	Data []CertificateInfo `json:"data"`
}
