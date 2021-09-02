// Package gmsm
//
// @author: xwc1125
// @date: 2020/3/2
package gmsm

import (
	"crypto/ecdsa"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

const (
	PubKeyBytesLenCompressed   = 33
	PubKeyBytesLenUncompressed = 65
)
const (
	pubkeyCompressed   byte = 0x2 // y_bit + x coord
	pubkeyUncompressed byte = 0x4 // x coord + y coord
)

func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

func MarshalPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	p := &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: Curve(),
			X:     key.X,
			Y:     key.Y,
		},
		D: key.D,
	}
	return x509.MarshalSm2UnecryptedPrivateKey(p)
}

func UnmarshalPrivateKey(keyBytes []byte) (*sm2.PrivateKey, error) {
	return x509.ParseSm2PrivateKey(keyBytes)
}
