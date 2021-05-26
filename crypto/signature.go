package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/crypto/gmsm"
	"github.com/chain5j/chain5j-pkg/crypto/prime256v1"
	"github.com/chain5j/chain5j-pkg/crypto/secp256k1"
	"github.com/chain5j/chain5j-pkg/json"
	"math/big"
)

type SignResult struct {
	Name      string `json:"name"`
	PubKey    []byte `json:"pubKey"`
	Signature []byte `json:"signature"`
}

type jsonSignResult struct {
	Name      string `json:"name"`
	PubKey    []byte `json:"pubKey"`
	Signature []byte `json:"signature"`
}

func (s *SignResult) Bytes() ([]byte, error) {
	return rlp.EncodeToBytes(&s)
}

func ParseSign(sig []byte) (*SignResult, error) {
	var result SignResult
	err := rlp.DecodeBytes(sig, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (s *SignResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(*s)
}

func (s *SignResult) UnmarshalJSON(bytes []byte) error {
	var result jsonSignResult
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	s.Signature = result.Signature
	s.PubKey = result.PubKey
	s.Name = result.Name
	return nil
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	result, err := ParseSign(sig)
	if err != nil {
		return nil, err
	}
	switch result.Name {
	case P256:
		return prime256v1.RecoverPubkey(hash, result.Signature)
	case S256:
		return secp256k1.RecoverPubkey(hash, result.Signature)
	case SM2P256:
		pubkey, err := gmsm.DecompressPubkey(result.PubKey)
		if err != nil {
			return nil, err
		}
		b := gmsm.VerifyECDSA(pubkey, hash, result.Signature)
		if b {
			return result.PubKey, nil
		}
		return nil, errors.New("SM2 verify is error")
	}

	return prime256v1.RecoverPubkey(hash, result.Signature)
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	sign, _ := ParseSign(sig)
	curve := CurveType(sign.Name)
	var (
		x, y *big.Int
	)
	switch sign.Name {
	case P256:
		x, y = elliptic.Unmarshal(curve, s)
	case S256:
		x, y = elliptic.Unmarshal(curve, s)
	case SM2P256:
		pubkey, err := UnmarshalPubkey(sign.Name, sign.PubKey)
		//pubkey, err := gmsm.DecompressPubkey(s)
		if err != nil {
			return nil, err
		}
		x = pubkey.X
		y = pubkey.Y
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

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
		sig1, err = prime256v1.SignCompact(prime256v1.NewPrivateKey(prv.D), hash, false)
	case S256:
		sig1, err = secp256k1.Sign(hash, prv)
	case SM2P256:
		sig1, err = gmsm.Sign(hash, prv)
		result.PubKey = gmsm.CompressPubkey(&prv.PublicKey)
		//result.PubKey = MarshalPubkey(&prv.PublicKey)
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
			return prime256v1.VerifySignature(pubkey, hash, signature[:len(signature)-1])
		} else {
			recoverPubkey, err := prime256v1.RecoverPubkey(hash, signature)
			if err != nil {
				return false
			}
			return prime256v1.VerifySignature(recoverPubkey, hash, signature[:len(signature)-1])
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
		unmarshalPubkey, err := UnmarshalPubkey(curveName, pubkey)
		if err != nil {
			return false
		}
		pubKey := gmsm.ECDSAToPubKey(unmarshalPubkey)
		return gmsm.Verify(pubKey, hash, signature)
	}
	return false
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(curveName string, pubkey []byte) (*ecdsa.PublicKey, error) {
	var x, y *big.Int
	switch curveName {
	case P256:
		x, y = prime256v1.DecompressPubkey(pubkey)
	case S256:
		return secp256k1.DecompressPubkey(pubkey)
	case SM2P256:
		return gmsm.DecompressPubkey(pubkey)
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
		return prime256v1.CompressPubkey(pubkey.X, pubkey.Y)
	case S256:
		return secp256k1.CompressPubkey(pubkey)
	case SM2P256:
		return gmsm.CompressPubkey(pubkey)
	}
	return nil
}
