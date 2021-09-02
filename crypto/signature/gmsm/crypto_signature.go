// Package gmsm
//
// @author: xwc1125
// @date: 2020/3/2
package gmsm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

// Sign 使用sm2进行签名
func Sign(hash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	privateKey := NewPrivateKey(prv.D)
	r, b, err := sm2.Sign(privateKey, hash)
	if err != nil {
		return nil, err
	}

	return sm2.SignDigitToSignData(r, b)
}

// VerifyECDSA 使用sm2 公钥验证签名
func VerifyECDSA(pub *ecdsa.PublicKey, msg []byte, signature []byte) bool {
	pubKey := FromECDSAPubKey(pub)
	return Verify(pubKey, msg, signature)
}

// Verify 使用sm2 公钥验证签名
func Verify(pub *sm2.PublicKey, msg []byte, signature []byte) bool {
	r, b, err := sm2.SignDataToSignDigit(signature)
	if err != nil {
		return false
	}
	return sm2.Verify(pub, msg, r, b)
}

// UnmarshalPublicKeyToECDSA 公钥反序列化
func UnmarshalPublicKeyToECDSA(pubkeyBytes []byte) (*ecdsa.PublicKey, error) {
	var x, y *big.Int
	publicKey := sm2.Decompress(pubkeyBytes)

	if publicKey == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	x = publicKey.X
	y = publicKey.Y
	return &ecdsa.PublicKey{X: x, Y: y, Curve: Curve()}, nil
}

// UnmarshalPublicKey 公钥反序列化
func UnmarshalPublicKey(pubkeyBytes []byte) (*sm2.PublicKey, error) {
	publicKey := sm2.Decompress(pubkeyBytes)

	if publicKey == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return publicKey, nil
}

// MarshalPublicKeyFromECDSA 公钥序列化
func MarshalPublicKeyFromECDSA(pubkey *ecdsa.PublicKey) []byte {
	p := &sm2.PublicKey{
		Curve: Curve(),
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	return sm2.Compress(p)
}

// MarshalPublicKey 公钥序列化
func MarshalPublicKey(pubkey *sm2.PublicKey) []byte {
	return sm2.Compress(pubkey)
}

// Curve 椭圆曲线
func Curve() elliptic.Curve {
	return sm2.P256Sm2()
}
