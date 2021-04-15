package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func DecryptToByte(apiv3Key, associatedData, nonce, ciphertext string) ([]byte, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher([]byte(apiv3Key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, []byte(nonce), decodedCiphertext, []byte(associatedData))
}

//  DecryptToString 将下载证书的回包解析成证书
//
//  解析后的证书是这样的
//  -----BEGIN CERTIFICATE-----
//	-----END CERTIFICATE-----
func DecryptToString(apiv3Key, associatedData, nonce, ciphertext string) (string, error) {
	certificateByte, err := DecryptToByte(apiv3Key, associatedData, nonce, ciphertext)
	return string(certificateByte), err
}
