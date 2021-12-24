// Package signature
//
// @author: xwc1125
package signature

import (
	"crypto"
	"hash"
)

// HashMsg hash data
func HashMsg(curveName string, msg []byte) []byte {
	ecdsa, err := GetECDSA(curveName)
	if err != nil {
		return nil
	}
	bytes, err := ecdsa.HashMsg(curveName, msg)
	if err != nil {
		return nil
	}
	return bytes
}

// HashType hash type
func HashType(curveName string) crypto.Hash {
	ecdsa, err := GetECDSA(curveName)
	if err != nil {
		return crypto.SHA256
	}
	return ecdsa.HashType(curveName)
}

// HashFunc get hash.Hash
func HashFunc(curveName string) func() hash.Hash {
	ecdsa, err := GetECDSA(curveName)
	if err != nil {
		return nil
	}
	return ecdsa.HashFunc(curveName)
}
