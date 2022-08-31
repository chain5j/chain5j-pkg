// Package secp256k1
//
// @author: xwc1125
package secp256k1

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"hash"

	"github.com/btcsuite/btcd/btcec"
	crypto2 "github.com/chain5j/chain5j-pkg/crypto"
	"github.com/chain5j/chain5j-pkg/crypto/hashalg/sha3"
)

var (
	_ crypto2.ECDSA = new(Secp251k1)
)

type Secp251k1 struct {
}

func (s Secp251k1) GenerateKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return GenerateKey()
}

func (s Secp251k1) HashType(curveName string) crypto.Hash {
	return crypto.SHA3_256
}

func (s Secp251k1) HashFunc(cryptoName string) func() hash.Hash {
	return sha3.NewKeccak256
}

func (s Secp251k1) HashMsg(cryptoName string, data []byte) ([]byte, error) {
	return sha3.Keccak256(data), nil
}

func (s Secp251k1) ToECDSA(prv crypto.PrivateKey) *ecdsa.PrivateKey {
	privateKey := prv.(*btcec.PrivateKey)
	return privateKey.ToECDSA()
}

func (s Secp251k1) FromECDSA(prv *ecdsa.PrivateKey) crypto.PrivateKey {
	return (*btcec.PrivateKey)(prv)
}

func (s Secp251k1) ToECDSAPubKey(pub crypto.PublicKey) *ecdsa.PublicKey {
	publicKey := pub.(*btcec.PublicKey)
	return publicKey.ToECDSA()
}

func (s Secp251k1) FromECDSAPubKey(pub *ecdsa.PublicKey) crypto.PublicKey {
	return (*btcec.PublicKey)(pub)
}

func (s Secp251k1) MarshalPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	return MarshalPrivateKey(key)
}

func (s Secp251k1) UnmarshalPrivateKey(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error) {
	return UnmarshalPrivateKey(curve, keyBytes)
}

func (s Secp251k1) MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	return MarshalPublicKey(pub)
}

func (s Secp251k1) UnmarshalPublicKey(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	return UnmarshalPublicKey(curve, data)
}

func (s Secp251k1) MarshalPrivateKeyX509(key *ecdsa.PrivateKey) ([]byte, error) {
	return MarshalPrivateKeyX509(key)
}

func (s Secp251k1) UnmarshalPrivateKeyX509(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error) {
	return UnmarshalPrivateKeyX509(curve, keyBytes)
}

func (s Secp251k1) MarshalPublicKeyX509(pub *ecdsa.PublicKey) ([]byte, error) {
	return MarshalPublicKeyX509(pub)
}

func (s Secp251k1) UnmarshalPublicKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	return UnmarshalPublicKeyX509(curve, data)
}

func (s Secp251k1) Sign(prv *ecdsa.PrivateKey, hash []byte) (sig []byte, err error) {
	return Sign(prv, hash)
}

func (s Secp251k1) Verify(pub *ecdsa.PublicKey, hash []byte, signature []byte) bool {
	return Verify(pub, hash, signature)
}

// GenerateKey 生成s256的私钥
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(S256(), rand.Reader)
}
