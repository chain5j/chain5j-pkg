// Package signature
//
// @author: xwc1125
package signature

import (
	"crypto/ecdsa"
	"errors"
)

// MarshalPubkeyWithECDSA marshal ecdsa publicKey to bytes
func MarshalPubkeyWithECDSA(pub *ecdsa.PublicKey) ([]byte, error) {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil, errors.New("pub is nil")
	}
	if ecdsa, err := GetECDSA(CurveName(pub.Curve)); err != nil {
		return nil, err
	} else {
		return ecdsa.MarshalPublicKey(pub)
	}
}

// UnmarshalPubkeyWithECDSA converts bytes to a  public key.
func UnmarshalPubkeyWithECDSA(curveName string, pub []byte) (*ecdsa.PublicKey, error) {
	if len(pub) == 0 {
		return nil, errors.New("pub bytes is empty")
	}

	if ecdsa, err := GetECDSA(curveName); err != nil {
		return nil, err
	} else {
		return ecdsa.UnmarshalPublicKey(CurveType(curveName), pub)
	}
}

func MarshalPubkeyWithECDSAX509(pub *ecdsa.PublicKey) ([]byte, error) {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil, errors.New("pub is nil")
	}
	if ecdsa, err := GetECDSA(CurveName(pub.Curve)); err != nil {
		return nil, err
	} else {
		return ecdsa.MarshalPublicKeyX509(pub)
	}
}

func MarshalPrvkeyWithECDSA(prv *ecdsa.PrivateKey) ([]byte, error) {
	if prv == nil {
		return nil, errors.New("ecdsa private key is empty")
	}

	if ecdsa, err := GetECDSA(CurveName(prv.Curve)); err != nil {
		return nil, err
	} else {
		return ecdsa.MarshalPrivateKey(prv)
	}
}

func UnMarshalPrvkeyWithECDSA(curveName string, key []byte) (*ecdsa.PrivateKey, error) {
	if key == nil {
		return nil, errors.New("private key is empty")
	}

	if ecdsa, err := GetECDSA(curveName); err != nil {
		return nil, err
	} else {
		return ecdsa.UnmarshalPrivateKey(CurveType(curveName), key)
	}
}

func MarshalPrvkeyWithECDSAX509(prv *ecdsa.PrivateKey) ([]byte, error) {
	if prv == nil {
		return nil, errors.New("ecdsa private key is empty")
	}
	if ecdsa, err := GetECDSA(CurveName(prv.Curve)); err != nil {
		return nil, err
	} else {
		return ecdsa.MarshalPrivateKeyX509(prv)
	}
}
