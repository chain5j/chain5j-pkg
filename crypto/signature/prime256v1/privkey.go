package prime256v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"io"
	"math/big"
)

const (
	// PrivKeyBytesLen defines the length in bytes of a serialized private key.
	PrivKeyBytesLen = 32
)

// PrivateKey wraps an ecdsa.PrivateKey as a convenience mainly for signing
// things with the the private key without having to directly import the ecdsa
// package.
type PrivateKey ecdsa.PrivateKey

func GenerateECDSAKeyWithRand(curve elliptic.Curve, rand io.Reader) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(curve, rand)
}

// GeneratePrivateKey is a wrapper for ecdsa.GenerateKey that returns a PrivateKey
// instead of the normal ecdsa.PrivateKey.
func GeneratePrivateKey(curve elliptic.Curve) (*PrivateKey, error) {
	key, err := GenerateKey(curve)
	if err != nil {
		return nil, err
	}
	return (*PrivateKey)(key), nil
}

// GetECDSAKeyValues get ecdsa privateKey values
func GetECDSAKeyValues(key *ecdsa.PrivateKey) (priv []byte, x, y *big.Int, err error) {
	priv = key.D.Bytes()
	x = key.PublicKey.X
	y = key.PublicKey.Y
	return
}

// NewPrivateKey instantiates a new private key from a scalar encoded as a
// big integer.
func NewPrivateKey(curve elliptic.Curve, d *big.Int) *PrivateKey {
	b := make([]byte, 0, PrivKeyBytesLen)
	dB := paddedAppend(PrivKeyBytesLen, b, d.Bytes())
	priv, _ := PrivKeyFromBytes(curve, dB)
	return priv
}

// PrivKeyFromBytes returns a private and public key for `curve' based on the
// private key passed as an argument as a byte slice.
func PrivKeyFromBytes(curve elliptic.Curve, pk []byte) (*PrivateKey, *PublicKey) {
	x, y := curve.ScalarBaseMult(pk)

	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(pk),
	}

	return (*PrivateKey)(priv), (*PublicKey)(&priv.PublicKey)
}

// Public returns the PublicKey corresponding to this private key.
func (p PrivateKey) Public() (*big.Int, *big.Int) {
	return p.PublicKey.X, p.PublicKey.Y
}

func (p PrivateKey) PublicBytes() (*big.Int, *big.Int) {
	return p.PublicKey.X, p.PublicKey.Y
}

// PubKey returns the PublicKey corresponding to this private key.
func (p *PrivateKey) PubKey() *PublicKey {
	return (*PublicKey)(&p.PublicKey)
}

// ToECDSA returns the private key as a *ecdsa.PrivateKey.
func (p *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return (*ecdsa.PrivateKey)(p)
}

// Sign generates an ECDSA signature for the provided hash (which should be the
// result of hashing a larger message) using the private key. Produced signature
// is deterministic (same message and same key yield the same signature) and
// canonical in accordance with RFC6979 and BIP0062.
func (p *PrivateKey) Sign(hash []byte) (*Signature, error) {
	return signRFC6979((*ecdsa.PrivateKey)(p), hash)
}

// Serialize returns the private key number d as a big-endian binary-encoded
// number, padded to a length of 32 bytes.
func (p PrivateKey) Serialize() []byte {
	b := make([]byte, 0, PrivKeyBytesLen)
	return paddedAppend(PrivKeyBytesLen, b, p.ToECDSA().D.Bytes())
}

// GetD satisfies the chainec PrivateKey interface.
func (p PrivateKey) GetD() *big.Int {
	return p.D
}
