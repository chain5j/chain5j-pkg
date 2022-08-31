// Package gmsm
//
// @author: xwc1125
package gmsm

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"hash"

	crypto2 "github.com/chain5j/chain5j-pkg/crypto"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
)

var (
	_ crypto2.ECDSA = new(SM2)
)

type SM2 struct {
}

func (s SM2) GenerateKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return GenerateKey()
}

func (s SM2) HashType(curveName string) crypto.Hash {
	return 20
}

func (s SM2) HashFunc(cryptoName string) func() hash.Hash {
	return sm3.New
}

func (s SM2) HashMsg(cryptoName string, data []byte) ([]byte, error) {
	return sm3.Sm3Sum(data), nil
}

func (s SM2) ToECDSA(prv crypto.PrivateKey) *ecdsa.PrivateKey {
	return ToECDSA(prv.(*sm2.PrivateKey))
}

func (s SM2) FromECDSA(prv *ecdsa.PrivateKey) crypto.PrivateKey {
	return FromECDSA(prv)
}

func (s SM2) ToECDSAPubKey(pub crypto.PublicKey) *ecdsa.PublicKey {
	return ToECDSAPubKey(pub.(*sm2.PublicKey))
}

func (s SM2) FromECDSAPubKey(pub *ecdsa.PublicKey) crypto.PublicKey {
	return FromECDSAPubKey(pub)
}

func (s SM2) MarshalPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	return MarshalPrivateKey(key)
}

func (s SM2) UnmarshalPrivateKey(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error) {
	privateKey, err := UnmarshalPrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return ToECDSA(privateKey), nil
}

func (s SM2) MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	return MarshalPublicKey(pub)
}

func (s SM2) UnmarshalPublicKey(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	return UnmarshalPublicKey(data)
}

func (s SM2) MarshalPrivateKeyX509(key *ecdsa.PrivateKey) ([]byte, error) {
	return MarshalPrivateKeyX509(key)
}

func (s SM2) UnmarshalPrivateKeyX509(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error) {
	return UnmarshalPrivateX509(keyBytes)
}

func (s SM2) MarshalPublicKeyX509(pub *ecdsa.PublicKey) ([]byte, error) {
	return MarshalPublicKeyX509(pub)
}

func (s SM2) UnmarshalPublicKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	return UnmarshalPublicKeyX509(data)
}

func (s SM2) Sign(prv *ecdsa.PrivateKey, msg []byte) (sig []byte, err error) {
	return Sign(prv, msg)
}

func (s SM2) Verify(pub *ecdsa.PublicKey, msg []byte, signature []byte) bool {
	return Verify(pub, msg, signature)
}

// GenerateKey 生成PrivateKey
func GenerateKey() (*ecdsa.PrivateKey, error) {
	key, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return ToECDSA(key), nil
}
