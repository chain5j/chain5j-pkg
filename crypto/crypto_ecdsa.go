// Package crypto
//
// @author: xwc1125
package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"hash"
)

type ECDSA interface {
	GenerateKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error)

	HashType(curveName string) crypto.Hash
	HashFunc(cryptoName string) func() hash.Hash
	HashMsg(cryptoName string, data []byte) ([]byte, error)

	ToECDSA(prv crypto.PrivateKey) *ecdsa.PrivateKey
	FromECDSA(prv *ecdsa.PrivateKey) crypto.PrivateKey
	ToECDSAPubKey(pub crypto.PublicKey) *ecdsa.PublicKey
	FromECDSAPubKey(pub *ecdsa.PublicKey) crypto.PublicKey

	MarshalPrivateKey(key *ecdsa.PrivateKey) ([]byte, error)
	UnmarshalPrivateKey(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error)
	MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error)
	UnmarshalPublicKey(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error)

	MarshalPrivateKeyX509(key *ecdsa.PrivateKey) ([]byte, error)
	UnmarshalPrivateKeyX509(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error)
	MarshalPublicKeyX509(pub *ecdsa.PublicKey) ([]byte, error)
	UnmarshalPublicKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error)

	Sign(prv *ecdsa.PrivateKey, hash []byte) (sig []byte, err error)
	Verify(pub *ecdsa.PublicKey, hash []byte, signature []byte) bool
}
