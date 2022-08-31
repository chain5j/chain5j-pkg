// Package gmsm
//
// @author: xwc1125
package gmsm

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/tjfoc/gmsm/sm2"
)

// Sign 使用sm2进行签名
func Sign(prv *ecdsa.PrivateKey, msg []byte) (sig []byte, err error) {
	privateKey := NewPrivateKey(prv.D)
	r, b, err := sm2.Sm2Sign(privateKey, msg, nil, rand.Reader)
	if err != nil {
		return nil, err
	}

	return sm2.SignDigitToSignData(r, b)
}

// Verify 使用sm2 公钥验证签名
func Verify(pub *ecdsa.PublicKey, msg []byte, signature []byte) bool {
	r, s, err := sm2.SignDataToSignDigit(signature)
	if err != nil {
		return false
	}
	return sm2.Sm2Verify(FromECDSAPubKey(pub), msg, nil, r, s)
}
