// Package secp256k1
//
// @author: xwc1125
package secp256k1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"

	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1/btcecv1"
)

func MarshalPrivateKey(key *ecdsa.PrivateKey) ([]byte, error) {
	privateKey := btcecv1.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: btcecv1.S256(),
			X:     key.X,
			Y:     key.Y,
		},
		D: key.D,
	}
	return privateKey.Serialize(), nil
}

func UnmarshalPrivateKey(curve elliptic.Curve, data []byte) (*ecdsa.PrivateKey, error) {
	privateKey, _ := btcecv1.PrivKeyFromBytes(btcecv1.S256(), data)
	if privateKey != nil {
		return privateKey.ToECDSA(), nil
	}
	return nil, errors.New("unmarshal private key err")
}

func MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	publicKey := btcecv1.PublicKey{
		Curve: pub.Curve,
		X:     pub.X,
		Y:     pub.Y,
	}
	return publicKey.SerializeUncompressed(), nil
}

func UnmarshalPublicKey(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	pubKey, err := btcecv1.ParsePubKey(data, btcecv1.S256())
	if err != nil {
		return nil, err
	}
	return pubKey.ToECDSA(), nil
}

// =================x509=================

func MarshalPrivateKeyX509(key *ecdsa.PrivateKey) ([]byte, error) {
	privateKey := btcecv1.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: btcecv1.S256(),
			X:     key.X,
			Y:     key.Y,
		},
		D: key.D,
	}
	return privateKey.Serialize(), nil
}

func UnmarshalPrivateKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PrivateKey, error) {
	privateKey, _ := btcecv1.PrivKeyFromBytes(btcecv1.S256(), data)
	if privateKey != nil {
		return privateKey.ToECDSA(), nil
	}
	return nil, errors.New("unmarshal private key err")
}

func MarshalPublicKeyX509(pub *ecdsa.PublicKey) ([]byte, error) {
	return MarshalPublicKey(pub)
}

func UnmarshalPublicKeyX509(curve elliptic.Curve, data []byte) (*ecdsa.PublicKey, error) {
	return UnmarshalPublicKey(curve, data)
}
