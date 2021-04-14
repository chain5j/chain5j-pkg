package gmsm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

func Sign(hash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	privateKey := NewPrivateKey(prv.D)
	r, b, err := sm2.Sign(privateKey, hash)
	if err != nil {
		return nil, err
	}

	return sm2.SignDigitToSignData(r, b)
}

func VerifyECDSA(pub *ecdsa.PublicKey, msg []byte, signature []byte) bool {
	pubKey := ECDSAToPubKey(pub)
	return Verify(pubKey, msg, signature)
}

func Verify(pub *sm2.PublicKey, msg []byte, signature []byte) bool {
	r, b, err := sm2.SignDataToSignDigit(signature)
	if err != nil {
		return false
	}
	return sm2.Verify(pub, msg, r, b)
}

func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	var x, y *big.Int
	publicKey := sm2.Decompress(pubkey)

	if publicKey == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	x = publicKey.X
	y = publicKey.Y
	return &ecdsa.PublicKey{X: x, Y: y, Curve: Curve256()}, nil
}

func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	p := &sm2.PublicKey{
		Curve: Curve256(),
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	return sm2.Compress(p)

}

func Curve256() elliptic.Curve {
	return sm2.P256Sm2()
}
