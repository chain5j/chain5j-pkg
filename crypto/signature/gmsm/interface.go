// Package gmsm
//
// @author: xwc1125
// @date: 2021/7/30
package gmsm

import "crypto/ecdsa"

type ECDSAPrivateKey interface {
	GenerateKey() (*ecdsa.PrivateKey, error)
}
