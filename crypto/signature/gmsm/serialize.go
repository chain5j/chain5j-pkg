// Package gmsm
//
// @author: xwc1125
package gmsm

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
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
	return x509.MarshalSm2PrivateKey(FromECDSA(key), nil)
}

func UnmarshalPrivateKey(keyBytes []byte) (*sm2.PrivateKey, error) {
	return x509.ParsePKCS8UnecryptedPrivateKey(keyBytes)
}

func MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	return sm2.Compress(FromECDSAPubKey(pub)), nil
}

func UnmarshalPublicKey(data []byte) (pub *ecdsa.PublicKey, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch p := r.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	publicKey := sm2.Decompress(data)
	if publicKey == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{X: publicKey.X, Y: publicKey.Y, Curve: Curve()}, nil
}

func MarshalPrivateX509(prv *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalSm2PrivateKey(FromECDSA(prv), nil)
}

func UnmarshalPrivateX509(data []byte) (*ecdsa.PrivateKey, error) {
	privateKey, err := x509.ParseSm2PrivateKey(data)
	if err != nil {
		return nil, err
	}
	return ToECDSA(privateKey), nil
}

// =================x509=================

func MarshalPrivateKeyX509(key *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalSm2UnecryptedPrivateKey(FromECDSA(key))
}

func UnmarshalPrivateKeyX509(keyBytes []byte) (*sm2.PrivateKey, error) {
	return x509.ParseSm2PrivateKey(keyBytes)
}

func MarshalPublicKeyX509(pub *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalSm2PublicKey(FromECDSAPubKey(pub))
}

func UnmarshalPublicKeyX509(data []byte) (*ecdsa.PublicKey, error) {
	publicKey, err := x509.ParseSm2PublicKey(data)
	if err != nil {
		return nil, err
	}
	return ToECDSAPubKey(publicKey), nil
}
