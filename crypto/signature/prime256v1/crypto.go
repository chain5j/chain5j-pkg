package prime256v1

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"hash"

	crypto2 "github.com/chain5j/chain5j-pkg/crypto"
)

// ecdsa: y^2 = x^3 + ax + b
// p256: y^2 = x^3 - 3x + b
// s256: y^2 = x^3 + b

var (
	_ crypto2.ECDSA = new(ECDSA)
)

type ECDSA struct {
}

func (e ECDSA) GenerateKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return GenerateKey(curve)
}

func (e ECDSA) HashType(curveName string) crypto.Hash {
	return HashType(curveName)
}

func (e ECDSA) HashFunc(cryptoName string) func() hash.Hash {
	return HashFunc(cryptoName)
}

func (e ECDSA) HashMsg(cryptoName string, data []byte) ([]byte, error) {
	return HashMsg(cryptoName, data), nil
}

func (e ECDSA) ToECDSA(prv crypto.PrivateKey) *ecdsa.PrivateKey {
	return prv.(*ecdsa.PrivateKey)
}

func (e ECDSA) FromECDSA(prv *ecdsa.PrivateKey) crypto.PrivateKey {
	return prv
}

func (e ECDSA) ToECDSAPubKey(pub crypto.PublicKey) *ecdsa.PublicKey {
	return pub.(*ecdsa.PublicKey)
}

func (e ECDSA) FromECDSAPubKey(pub *ecdsa.PublicKey) crypto.PublicKey {
	return pub
}

func (e ECDSA) MarshalPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	return MarshalPrivateKey(key)
}

func (e ECDSA) UnmarshalPrivateKey(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error) {
	return UnmarshalPrivateKey(curve, keyBytes)
}

func (e ECDSA) MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	return MarshalPublicKey(pub)
}

func (e ECDSA) UnmarshalPublicKey(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	return UnmarshalPublicKey(curve, data)
}

func (e ECDSA) MarshalPrivateKeyX509(key *ecdsa.PrivateKey) ([]byte, error) {
	return MarshalPrivateKeyX509(key)
}

func (e ECDSA) UnmarshalPrivateKeyX509(curve elliptic.Curve, keyBytes []byte) (*ecdsa.PrivateKey, error) {
	return UnmarshalPrivateKeyX509(curve, keyBytes)
}

func (e ECDSA) MarshalPublicKeyX509(pub *ecdsa.PublicKey) ([]byte, error) {
	return MarshalPublicKeyX509(pub)
}

func (e ECDSA) UnmarshalPublicKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	return UnmarshalPublicKeyX509(curve, data)
}

func (e ECDSA) Sign(prv *ecdsa.PrivateKey, hash []byte) (sig []byte, err error) {
	return Sign(prv, hash)
}

func (e ECDSA) Verify(pub *ecdsa.PublicKey, hash []byte, signature []byte) bool {
	return Verify(pub, hash, signature)
}

// GenerateKey 生成私钥
func GenerateKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	return GenerateECDSAKeyWithRand(curve, rand.Reader)
}
