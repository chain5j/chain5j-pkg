// Package crypto
//
// @author: xwc1125
// @date: 2021/8/10
package crypto

import (
	"crypto"
	"crypto/sha256"
	"crypto/sha512"
	sha32 "github.com/chain5j/chain5j-pkg/crypto/hash/sha3"
	"github.com/chain5j/chain5j-pkg/crypto/signature"
	"github.com/tjfoc/gmsm/sm3"
	"hash"
)

// Sha256 sha256
func Sha256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// Sha384 Sha384
func Sha384(data []byte) []byte {
	h := sha512.New384()
	h.Write(data)
	return h.Sum(nil)
}

// Sha512 Sha512
func Sha512(data []byte) []byte {
	h := sha512.New()
	h.Write(data)
	return h.Sum(nil)
}

// HashMsg hash data
func HashMsg(curveName string, msg []byte) []byte {
	switch curveName {
	case signature.P256:
		return Sha256(msg)
	case signature.P384:
		return Sha384(msg)
	case signature.P521:
		return Sha512(msg)
	case signature.S256:
		return sha32.Keccak256(msg)
	case signature.SM2P256:
		return sm3.Sm3Sum(msg)
	default:
		return nil
	}
}

// HashType hash type
func HashType(curveName string) crypto.Hash {
	switch curveName {
	case signature.P256:
		return crypto.SHA256
	case signature.P384:
		return crypto.SHA384
	case signature.P521:
		return crypto.SHA512
	case signature.S256:
		return crypto.SHA3_256
	case signature.SM2P256:
		return 20
	default:
		return crypto.SHA256
	}
}

// HashFunc get hash.Hash
func HashFunc(cryptoName string) func() hash.Hash {
	switch cryptoName {
	case signature.P256:
		return sha256.New
	case signature.P384:
		return sha512.New384
	case signature.P521:
		return sha512.New
	case signature.S256:
		return sha32.NewKeccak256
	case signature.SM2P256:
		return sm3.New
	default:
		return nil
	}
}
