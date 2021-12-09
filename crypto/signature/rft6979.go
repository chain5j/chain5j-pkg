// Package signature
//
// @author: xwc1125
package signature

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"errors"
	"hash"
	"math/big"
)

var (
	// Used in RFC6979 implementation when testing the nonce for correctness
	one = big.NewInt(1)

	// oneInitializer is used to fill a byte slice with byte 0x01.  It is provided
	// here to avoid the need to create it multiple times.
	oneInitializer = []byte{0x01}
)

// NewSignature instantiates a new signature given some R,S values.
func NewSignature(curve elliptic.Curve, r, s *big.Int) *SignatureECDSA {
	order := new(big.Int).Set(curve.Params().N)
	return &SignatureECDSA{
		curve:     curve,
		order:     order,
		halfOrder: new(big.Int).Rsh(order, 1),
		R:         r,
		S:         s,
	}
}

// SignatureECDSA is a type representing an ecdsa signature.
type SignatureECDSA struct {
	curve elliptic.Curve
	// Curve order and halfOrder, used to tame ECDSA malleability (see BIP-0062)
	order     *big.Int
	halfOrder *big.Int
	R         *big.Int
	S         *big.Int
}

func (sig *SignatureECDSA) Verify(hash []byte, pubKey *ecdsa.PublicKey) bool {
	return ecdsa.Verify(pubKey, hash, sig.GetR(), sig.GetS())
}

// GetR satisfies the chainec PublicKey interface.
func (sig *SignatureECDSA) GetR() *big.Int {
	return sig.R
}

// GetS satisfies the chainec PublicKey interface.
func (sig *SignatureECDSA) GetS() *big.Int {
	return sig.S
}

func getOrder(curve elliptic.Curve) *big.Int {
	return new(big.Int).Set(curve.Params().N)
}

func getHalforder(curve elliptic.Curve) *big.Int {
	return new(big.Int).Rsh(getOrder(curve), 1)
}

// SignRFC6979 produces a compact signature of the data in hash with the given
// private key on the given koblitz curve. The isCompressed  parameter should
// be used to detail if the given signature should reference a compressed
// public key or not. If successful the bytes of the compact signature will be
// returned in the format:
// <(byte of 27+public key solution)+4 if compressed >< padded bytes for signature R><padded bytes for signature S>
// where the R and S parameters are padde up to the bitlengh of the curve.
func SignRFC6979(key *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	curve := key.Curve
	sig, err := signRFC6979(key, hash)
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, 2*curve.Params().BitSize)
	curvelen := (curve.Params().BitSize + 7) / 8

	// Pad R and S to curvelen if needed.
	bytelen := (sig.R.BitLen() + 7) / 8
	if bytelen < curvelen {
		result = append(result,
			make([]byte, curvelen-bytelen)...)
	}
	result = append(result, sig.R.Bytes()...)
	bytelen = (sig.S.BitLen() + 7) / 8
	if bytelen < curvelen {
		result = append(result,
			make([]byte, curvelen-bytelen)...)
	}
	result = append(result, sig.S.Bytes()...)
	return result, nil
}

// signRFC6979 generates a deterministic ECDSA signature according to RFC 6979
// and BIP 62.
func signRFC6979(privkey *ecdsa.PrivateKey, hash []byte) (*SignatureECDSA, error) {
	curve := privkey.Curve
	N := new(big.Int).Set(curve.Params().N)
	// 获取唯一的nonce
	k := NonceRFC6979(curve, privkey.D, hash, nil, nil)

	inv := new(big.Int).ModInverse(k, N)
	r, _ := curve.ScalarBaseMult(k.Bytes())
	r.Mod(r, N)

	if r.Sign() == 0 {
		return nil, errors.New("calculated R is zero")
	}

	e := hashToInt(curve, hash)
	s := new(big.Int).Mul(privkey.D, r)
	s.Add(s, e)
	s.Mul(s, inv)
	s.Mod(s, N)

	if s.Cmp(getHalforder(curve)) == 1 {
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		return nil, errors.New("calculated S is zero")
	}
	return &SignatureECDSA{
		curve: curve,
		R:     r,
		S:     s,
	}, nil
}

// NonceRFC6979 generates an ECDSA nonce (`k`) deterministically according to
// RFC 6979. It takes a 32-byte hash as an input and returns 32-byte nonce to
// be used in ECDSA algorithm.
func NonceRFC6979(curve elliptic.Curve, privkey *big.Int, hash []byte, extra []byte, version []byte) *big.Int {
	q := curve.Params().N
	x := privkey
	alg := HashFunc(curve.Params().Name)

	qlen := q.BitLen()
	holen := alg().Size()
	rolen := (qlen + 7) >> 3
	bx := append(int2octets(x, rolen), bits2octets(curve, hash, rolen)...)
	if len(extra) == 32 {
		bx = append(bx, extra...)
	}
	if len(version) == 16 && len(extra) == 32 {
		bx = append(bx, extra...)
	}
	if len(version) == 16 && len(extra) != 32 {
		bx = append(bx, bytes.Repeat([]byte{0x00}, 32)...)
		bx = append(bx, version...)
	}

	// Step B
	v := bytes.Repeat(oneInitializer, holen)

	// Step C (Go zeroes the all allocated memory)
	k := make([]byte, holen)

	// Step D
	k = mac(alg, k, append(append(v, 0x00), bx...))

	// Step E
	v = mac(alg, k, v)

	// Step F
	k = mac(alg, k, append(append(v, 0x01), bx...))

	// Step G
	v = mac(alg, k, v)

	// Step H
	for {
		// Step H1
		var t []byte

		// Step H2
		for len(t)*8 < qlen {
			v = mac(alg, k, v)
			t = append(t, v...)
		}

		// Step H3
		secret := hashToInt(curve, t)
		if secret.Cmp(one) >= 0 && secret.Cmp(q) < 0 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00))
		v = mac(alg, k, v)
	}
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// left pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(curve elliptic.Curve, in []byte, rolen int) []byte {
	z1 := hashToInt(curve, in)
	z2 := new(big.Int).Sub(z1, curve.Params().N)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(curve elliptic.Curve, hash []byte) *big.Int {
	orderBits := curve.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

// VerifyRFC6979 验证签名
func VerifyRFC6979(pubkey *ecdsa.PublicKey, msg, signature []byte) bool {
	bitlen := (pubkey.Curve.Params().BitSize + 7) / 8
	if len(signature) != bitlen*2 {
		return false
	}

	// 校验码在最后一位
	sig := &SignatureECDSA{
		curve: pubkey.Curve,
		R:     new(big.Int).SetBytes(signature[:bitlen]),
		S:     new(big.Int).SetBytes(signature[bitlen : bitlen*2]),
	}
	verify := sig.Verify(msg, pubkey)
	return verify
}
