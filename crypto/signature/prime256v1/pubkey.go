package prime256v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"
)

// These constants define the lengths of serialized public keys.
const (
	PubKeyBytesLenCompressed   = 33
	PubKeyBytesLenUncompressed = 65

	pubkeyCompressed   byte = 0x2 // y_bit + x coord
	pubkeyUncompressed byte = 0x4 // x coord + y coord
)

func isOdd(a *big.Int) bool {
	return a.Bit(0) == 1
}

// decompressPoint decompresses a point on the given curve given the X point and
// the solution to use.
func decompressPoint(curve elliptic.Curve, x *big.Int, ybit bool) (*big.Int, error) {
	// Y = +-sqrt(x^3 - 3x + B)
	x3 := new(big.Int).Mul(x, x)
	xThree := new(big.Int).Mul(x, big.NewInt(int64(3)))
	x3.Mul(x3, x)

	x3.Sub(x3, xThree)
	x3.Add(x3, curve.Params().B)

	var q = new(big.Int).Div(new(big.Int).Add(curve.Params().P,
		big.NewInt(1)), big.NewInt(4))
	// now calculate sqrt mod p of x2 + B
	// This code used to do a full sqrt based on tonelli/shanks,
	// but this was replaced by the algorithms referenced in
	// https://bitcointalk.org/index.php?topic=162805.msg1712294#msg1712294
	y := new(big.Int).Exp(x3, q, curve.Params().P)

	if ybit != isOdd(y) {
		y.Sub(curve.Params().P, y)
	}
	if ybit != isOdd(y) {
		return nil, fmt.Errorf("ybit doesn't match oddness")
	}
	return y, nil
}

// NewPublicKey instantiates a new public key with the given X,Y coordinates.
func NewPublicKey(curve elliptic.Curve, x *big.Int, y *big.Int) *PublicKey {
	return &PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
}

// ParsePubKey parses a public key for curve from a bytestring into a
// ecdsa.Publickey, verifying that it is valid. It supports compressed and
// uncompressed signature formats, but not the hybrid format.
func ParsePubKey(curve elliptic.Curve, pubKeyStr []byte) (key *PublicKey, err error) {
	pubkey := PublicKey{
		Curve: curve,
	}

	if len(pubKeyStr) == 0 {
		return nil, errors.New("pubkey string is empty")
	}

	format := pubKeyStr[0]
	ybit := (format & 0x1) == 0x1
	format &= ^byte(0x1)

	switch len(pubKeyStr) {
	case PubKeyBytesLenUncompressed:
		if format != pubkeyUncompressed {
			return nil, fmt.Errorf("invalid magic in pubkey str: "+
				"%d", pubKeyStr[0])
		}

		pubkey.X = new(big.Int).SetBytes(pubKeyStr[1:33])
		pubkey.Y = new(big.Int).SetBytes(pubKeyStr[33:])
	case PubKeyBytesLenCompressed:
		// format is 0x2 | solution, <X coordinate>
		// solution determines which solution of the curve we use.
		// / y^2 = x^3 + a*x + Curve.B
		if format != pubkeyCompressed {
			return nil, fmt.Errorf("invalid magic in compressed "+
				"pubkey string: %d", pubKeyStr[0])
		}
		pubkey.X = new(big.Int).SetBytes(pubKeyStr[1:33])
		pubkey.Y, err = decompressPoint(curve, pubkey.X, ybit)
		if err != nil {
			return nil, err
		}
	default: // wrong!
		return nil, fmt.Errorf("invalid pub key length %d",
			len(pubKeyStr))
	}

	if pubkey.X.Cmp(pubkey.Curve.Params().P) >= 0 {
		return nil, fmt.Errorf("pubkey X parameter is >= to P")
	}
	if pubkey.Y.Cmp(pubkey.Curve.Params().P) >= 0 {
		return nil, fmt.Errorf("pubkey Y parameter is >= to P")
	}
	if !pubkey.Curve.IsOnCurve(pubkey.X, pubkey.Y) {
		return nil, fmt.Errorf("pubkey [%v,%v] isn't on secp256k1 curve",
			pubkey.X, pubkey.Y)
	}
	return &pubkey, nil
}

// PublicKey is an ecdsa.PublicKey with additional functions to
// serialize in uncompressed and compressed formats.
type PublicKey ecdsa.PublicKey

// ToECDSA returns the public key as a *ecdsa.PublicKey.
func (p PublicKey) ToECDSA() *ecdsa.PublicKey {
	ecpk := ecdsa.PublicKey(p)
	return &ecpk
}

// Serialize serializes a public key in a 33-byte compressed format.
// It is the default serialization method.
func (p PublicKey) Serialize() []byte {
	return p.SerializeCompressed()
}

// SerializeCompressed serializes a public key in a 33-byte compressed format.
func (p PublicKey) SerializeCompressed() []byte {
	if p.Curve != elliptic.P256() {
		return nil
	}
	b := make([]byte, 0, PubKeyBytesLenCompressed)
	format := pubkeyCompressed
	if isOdd(p.Y) {
		format |= 0x1
	}
	b = append(b, format)
	return paddedAppend(32, b, p.X.Bytes())
}

// SerializeUncompressed serializes a public key in a 65-byte uncompressed
// format.
func (p PublicKey) SerializeUncompressed() []byte {
	b := make([]byte, 0, PubKeyBytesLenUncompressed)
	b = append(b, pubkeyUncompressed)
	b = paddedAppend(32, b, p.X.Bytes())
	return paddedAppend(32, b, p.Y.Bytes())
}

// IsEqual compares this PublicKey instance to the one passed, returning true if
// both PublicKeys are equivalent. A PublicKey is equivalent to another, if they
// both have the same X and Y coordinate.
func (p *PublicKey) IsEqual(otherPubKey *PublicKey) bool {
	return p.X.Cmp(otherPubKey.X) == 0 &&
		p.Y.Cmp(otherPubKey.Y) == 0
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// GetCurve satisfies the chainec PublicKey interface.
func (p PublicKey) GetCurve() elliptic.Curve {
	return p.Curve
}

// GetX satisfies the chainec PublicKey interface.
func (p PublicKey) GetX() *big.Int {
	return p.X
}

// GetY satisfies the chainec PublicKey interface.
func (p PublicKey) GetY() *big.Int {
	return p.Y
}

func DecompressPubkey(curve elliptic.Curve, pubkey []byte) (x, y *big.Int) {
	if curve != elliptic.P256() {
		return nil, nil
	}
	key, err := ParsePubKey(curve, pubkey)
	if err != nil {
		return nil, nil
	}
	return key.X, key.Y
}

func CompressPubkey(curve elliptic.Curve, x, y *big.Int) []byte {
	publicKey := NewPublicKey(curve, x, y)
	return publicKey.SerializeCompressed()
}
