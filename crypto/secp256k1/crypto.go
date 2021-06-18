// Package secp256k1
// 
// @author: xwc1125
// @date: 2021/6/18
package secp256k1

import (
	"crypto/ecdsa"
	"crypto/rand"
)

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(S256(), rand.Reader)
}
