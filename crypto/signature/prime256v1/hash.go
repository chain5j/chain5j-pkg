// Package prime256v1
//
// @author: xwc1125
package prime256v1

import (
	"crypto"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// HashFunc get hash.Hash
func HashFunc(cryptoName string) func() hash.Hash {
	switch cryptoName {
	case "P-256":
		return sha256.New
	case "P-384":
		return sha512.New384
	case "P-521":
		return sha512.New
	default:
		return nil
	}
}

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
	case "P-256":
		return Sha256(msg)
	case "P-384":
		return Sha384(msg)
	case "P-521":
		return Sha512(msg)
	default:
		return nil
	}
}

// HashType hash type
func HashType(curveName string) crypto.Hash {
	switch curveName {
	case "P-256":
		return crypto.SHA256
	case "P-384":
		return crypto.SHA384
	case "P-521":
		return crypto.SHA512
	default:
		return crypto.SHA256
	}
}
