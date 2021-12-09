// Package prime256v1
//
// @author: xwc1125
package prime256v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/chain5j/chain5j-pkg/math"
	"math/big"
)

func MarshalPrivateKey(prv *ecdsa.PrivateKey) ([]byte, error) {
	if prv.Curve == elliptic.P256() {
		return p256MarshalPrivateKey(prv)
	}
	return MarshalPrivateKeyX509(prv)
}

func p256MarshalPrivateKey(prv *ecdsa.PrivateKey) ([]byte, error) {
	return math.PaddedBigBytes(prv.D, prv.Params().BitSize/8), nil
}

func UnmarshalPrivateKey(curve elliptic.Curve, data []byte) (*ecdsa.PrivateKey, error) {
	if curve == elliptic.P256() {
		return toECDSA(curve, data, true)
	}
	return UnmarshalPrivateKeyX509(curve, data)
}

func toECDSA(curve elliptic.Curve, d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)

	priv.PublicKey.Curve = curve
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(priv.PublicKey.Curve.Params().N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

func MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	if pub.Curve == elliptic.P256() {
		return CompressPubkey(pub.Curve, pub.X, pub.Y), nil
	} else {
		return MarshalPublicKeyX509(pub)
	}
}

func UnmarshalPublicKey(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	if curve == elliptic.P256() {
		x, y := DecompressPubkey(curve, data)
		return &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}, nil
	} else {
		return UnmarshalPublicKeyX509(curve, data)
	}
}

// =================x509=================

func MarshalPrivateKeyX509(prv *ecdsa.PrivateKey) ([]byte, error) {
	return x509.MarshalECPrivateKey(prv)
}

func UnmarshalPrivateKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PrivateKey, error) {
	return x509.ParseECPrivateKey(data)
}

func MarshalPublicKeyX509(pub *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

func UnmarshalPublicKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	publicKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}
	return publicKey.(*ecdsa.PublicKey), nil
}
