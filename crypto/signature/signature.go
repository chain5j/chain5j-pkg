// Package signature
//
// @author: xwc1125
package signature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"github.com/chain5j/chain5j-pkg/crypto/signature/gmsm"
	"github.com/chain5j/chain5j-pkg/crypto/signature/prime256v1"
	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1"
	"math/big"
)

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash []byte, sig *SignResult) ([]byte, error) {
	if sig == nil {
		return nil, errors.New("sign result is empty")
	}
	switch sig.Name {
	case P256:
		return prime256v1.RecoverPubkey(elliptic.P256(), hash, sig.Signature)
	case S256:
		return secp256k1.RecoverPubkey(hash, sig.Signature)
	case SM2P256:
		pubkey, err := gmsm.UnmarshalPublicKey(sig.PubKey)
		if err != nil {
			return nil, err
		}
		b := gmsm.Verify(pubkey, hash, sig.Signature)
		if b {
			return sig.PubKey, nil
		}
		return nil, errors.New("SM2 verify is error")
	default:
		return nil, errors.New("unsupported signName")
	}
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash []byte, sig *SignResult) (*ecdsa.PublicKey, error) {
	if sig == nil {
		return nil, errors.New("sign result is empty")
	}
	if sig.PubKey.Nil() {
		s, err := Ecrecover(hash, sig)
		if err != nil {
			return nil, err
		}
		sig.PubKey = s
	}
	return UnmarshalPubkeyWithECDSA(sig.Name, sig.PubKey)
}

// BigintToPub 将big.Int转换为PublicKey
func BigintToPub(curveName string, x, y *big.Int) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{Curve: CurveType(curveName), X: x, Y: y}
}
