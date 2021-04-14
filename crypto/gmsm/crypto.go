// description: chain5j-core
// 
// @author: xwc1125
// @date: 2020/3/2
package gmsm

import (
	"crypto/ecdsa"
	"errors"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

var errInvalidPubkey = errors.New("invalid public key")

func GenerateKey() (*ecdsa.PrivateKey, error) {
	key, err := sm2.GenerateKey()
	if err != nil {
		return nil, err
	}
	return PrivKeyToECDSA(key), nil
}

func PrivKeyToECDSA(p *sm2.PrivateKey) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.D = p.D
	priv.PublicKey = *PubKeyToECDSA(&p.PublicKey)
	return priv
}

func ECDSAToPrivKey(p *ecdsa.PrivateKey) *sm2.PrivateKey {
	priv := new(sm2.PrivateKey)
	priv.D = p.D
	priv.PublicKey = *ECDSAToPubKey(&p.PublicKey)
	return priv
}

func PubKeyToECDSA(pub *sm2.PublicKey) *ecdsa.PublicKey {
	return (*ecdsa.PublicKey)(pub)
}

func ECDSAToPubKey(pub *ecdsa.PublicKey) *sm2.PublicKey {
	return (*sm2.PublicKey)(pub)
}

const PrivKeyBytesLen = 32

func NewPrivateKey(d *big.Int) *sm2.PrivateKey {
	b := make([]byte, 0, PrivKeyBytesLen)
	dB := paddedAppend(PrivKeyBytesLen, b, d.Bytes())
	priv, _ := PrivKeyFromBytes(dB)
	return priv
}

func PrivKeyFromBytes(pk []byte) (*sm2.PrivateKey, *sm2.PublicKey) {
	curve := Curve256()
	x, y := curve.ScalarBaseMult(pk)

	priv := &sm2.PrivateKey{
		PublicKey: sm2.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return priv, &priv.PublicKey
}
