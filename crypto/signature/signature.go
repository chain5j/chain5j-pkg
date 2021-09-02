// Package signature
//
// @author: xwc1125
// @date: 2019/9/9
package signature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/crypto/signature/gmsm"
	"github.com/chain5j/chain5j-pkg/crypto/signature/prime256v1"
	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1"
	"math/big"
)

// SignResult 签名结果
type SignResult struct {
	Name      string `json:"name" mapstructure:"name"`           // 算法名称
	PubKey    []byte `json:"pub_key" mapstructure:"pub_key"`     // 公钥
	Signature []byte `json:"signature" mapstructure:"signature"` // 签名结果
}

// Bytes 获取signResult的rlp bytes
func (s *SignResult) Bytes() ([]byte, error) {
	return rlp.EncodeToBytes(&s)
}

// ParseSign rlp bytes to signResult
func ParseSign(sig []byte) (*SignResult, error) {
	var result SignResult
	err := rlp.DecodeBytes(sig, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	result, err := ParseSign(sig)
	if err != nil {
		return nil, err
	}
	switch result.Name {
	case P256:
		return prime256v1.RecoverPubkey(elliptic.P256(), hash, result.Signature)
	case S256:
		return secp256k1.RecoverPubkey(hash, result.Signature)
	case SM2P256:
		pubkey, err := gmsm.UnmarshalPublicKeyToECDSA(result.PubKey)
		if err != nil {
			return nil, err
		}
		b := gmsm.VerifyECDSA(pubkey, hash, result.Signature)
		if b {
			return result.PubKey, nil
		}
		return nil, errors.New("SM2 verify is error")
	default:
		return nil, errors.New("unsupported signName")
	}
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	sign, err := ParseSign(sig)
	if err != nil {
		return nil, err
	}
	if len(sign.PubKey) == 0 {
		s, err := Ecrecover(hash, sig)
		if err != nil {
			return nil, err
		}
		sign.PubKey = s
	}
	return UnmarshalPubkeyWithECDSA(sign.Name, sign.PubKey)
}

// BigintToPub TODO
func BigintToPub(curveName string, x, y *big.Int) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{Curve: CurveType(curveName), X: x, Y: y}
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given hash cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(hash []byte, prv *ecdsa.PrivateKey) (sig *SignResult, err error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	signType := CurveName(prv.Curve)
	result := SignResult{
		Name: signType,
	}

	var (
		sig1 []byte
	)
	switch signType {
	case P256:
		sig1, err = prime256v1.SignCompact(prime256v1.NewPrivateKey(elliptic.P256(), prv.D), hash, false)
	case S256:
		sig1, err = secp256k1.Sign(hash, prv)
	case SM2P256:
		sig1, err = gmsm.Sign(hash, prv)
		result.PubKey = gmsm.MarshalPublicKeyFromECDSA(&prv.PublicKey)
	default:
		return nil, errors.New("unsupported signType")
	}
	if err != nil {
		return nil, err
	}
	result.Signature = sig1
	return &result, nil
}

// VerifySignature checks that the given public key created signature over hash.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature(curveName string, pubkey, hash, signature []byte) bool {
	switch curveName {
	case P256:
		if pubkey != nil && len(pubkey) > 0 {
			publicKey, err := UnmarshalPubkeyWithECDSA(curveName, pubkey)
			if err != nil {
				return false
			}
			return prime256v1.VerifySignature(publicKey, hash, signature[:len(signature)-1])
		} else {
			recoverPubkey, _, err := prime256v1.RecoverCompact(elliptic.P256(), hash, signature)
			if err != nil {
				return false
			}
			return prime256v1.VerifySignature(recoverPubkey.ToECDSA(), hash, signature[:len(signature)-1])
		}
	case S256:
		if pubkey != nil && len(pubkey) > 0 {
			return secp256k1.VerifySignature(pubkey, hash, signature[:len(signature)-1])
		} else {
			recoverPubkey, err := secp256k1.RecoverPubkey(hash, signature)
			if err != nil {
				return false
			}
			return secp256k1.VerifySignature(recoverPubkey, hash, signature[:len(signature)-1])
		}
	case SM2P256:
		unmarshalPubkey, err := UnmarshalPubkeyWithECDSA(curveName, pubkey)
		if err != nil {
			return false
		}
		pubKey := gmsm.FromECDSAPubKey(unmarshalPubkey)
		return gmsm.Verify(pubKey, hash, signature)
	}
	return false
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(curveName string, pubkey []byte) (*ecdsa.PublicKey, error) {
	var x, y *big.Int
	switch curveName {
	case P256:
		x, y = prime256v1.DecompressPubkey(elliptic.P256(), pubkey)
	case S256:
		return secp256k1.DecompressPubkey(pubkey)
	case SM2P256:
		return gmsm.UnmarshalPublicKeyToECDSA(pubkey)
	}

	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{X: x, Y: y, Curve: CurveType(curveName)}, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	curveName := CurveName(pubkey.Curve)
	switch curveName {
	case P256:
		return prime256v1.CompressPubkey(elliptic.P256(), pubkey.X, pubkey.Y)
	case S256:
		return secp256k1.CompressPubkey(pubkey)
	case SM2P256:
		return gmsm.MarshalPublicKeyFromECDSA(pubkey)
	}
	return nil
}
