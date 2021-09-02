// Package prime256v1
//
// @author: xwc1125
// @date: 2021/8/19
package prime256v1

import (
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
