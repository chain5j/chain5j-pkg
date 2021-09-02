// Package gmsm
//
// @author: xwc1125
// @date: 2020/3/2
package gmsm

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

var errInvalidPubkey = errors.New("invalid public key")

// GenerateKey 生成PrivateKey
func GenerateKey() (*ecdsa.PrivateKey, error) {
	key, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return ToECDSA(key), nil
}

// ToECDSA sm2.PrivateKey to ecdsa.PrivateKey
func ToECDSA(p *sm2.PrivateKey) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.D = p.D
	priv.PublicKey = *ToECDSAPubKey(&p.PublicKey)
	return priv
}

// FromECDSA ecdsa.PrivateKey to sm2.PrivateKey
func FromECDSA(p *ecdsa.PrivateKey) *sm2.PrivateKey {
	priv := new(sm2.PrivateKey)
	priv.D = p.D
	priv.PublicKey = *FromECDSAPubKey(&p.PublicKey)
	return priv
}

// ToECDSAPubKey sm2.PublicKey to ecdsa.PublicKey
func ToECDSAPubKey(pub *sm2.PublicKey) *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(pub)
}

// FromECDSAPubKey ecdsa.PublicKey to sm2.PublicKey
func FromECDSAPubKey(pub *ecdsa.PublicKey) *sm2.PublicKey {
	return (*sm2.PublicKey)(pub)
}

// PrivKeyBytesLen 私钥长度
const PrivKeyBytesLen = 32

// NewPrivateKey new sm2.PrivateKey
func NewPrivateKey(d *big.Int) *sm2.PrivateKey {
	b := make([]byte, 0, PrivKeyBytesLen)
	dB := paddedAppend(PrivKeyBytesLen, b, d.Bytes())
	priv, _ := PrivKeyFromBytes(dB)
	return priv
}

// PrivKeyFromBytes bytes to sm2.privateKey
func PrivKeyFromBytes(pk []byte) (*sm2.PrivateKey, *sm2.PublicKey) {
	curve := Curve()
	x, y := curve.ScalarBaseMult(pk)

	priv := &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return priv, &priv.PublicKey
}
