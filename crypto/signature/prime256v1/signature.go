package prime256v1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var (
	errNegativeValue          = errors.New("value may be interpreted as negative")
	errExcessivelyPaddedValue = errors.New("value is excessively padded")
)

func Sign(prv *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	if prv.Curve == elliptic.P256() {
		return SignCompact((*PrivateKey)(prv), hash, true)
	}
	return ecdsa.SignASN1(rand.Reader, prv, hash)
}

func Verify(pub *ecdsa.PublicKey, hash []byte, signature []byte) bool {
	if pub.Curve == elliptic.P256() {
		return VerifySignature(pub, hash, signature)
	}
	return ecdsa.VerifyASN1(pub, hash, signature)
}

// Serialize returns the ECDSA signature in the more strict DER format.  Note
// that the serialized bytes returned do not include the appended hash type
// used in Decred signature scripts.
//
// encoding/asn1 is broken so we hand roll this output:
//
// 0x30 <length> 0x02 <length r> r 0x02 <length s> s
func (sig *Signature) Serialize() []byte {
	// low 'S' malleability breaker
	sigS := sig.S
	if sigS.Cmp(getOrder(sig.curve)) == 1 {
		sigS = new(big.Int).Sub(getOrder(sig.curve), sigS)
	}
	// Ensure the encoded bytes for the r and s values are canonical and
	// thus suitable for DER encoding.
	rb := canonicalizeInt(sig.R)
	sb := canonicalizeInt(sigS)

	// total length of returned signature is 1 byte for each magic and
	// length (6 total), plus lengths of r and s
	length := 6 + len(rb) + len(sb)
	b := make([]byte, length)

	b[0] = 0x30
	b[1] = byte(length - 2)
	b[2] = 0x02
	b[3] = byte(len(rb))
	offset := copy(b[4:], rb) + 4
	b[offset] = 0x02
	b[offset+1] = byte(len(sb))
	copy(b[offset+2:], sb)
	return b
}

// IsEqual compares this Signature instance to the one passed, returning true
// if both Signatures are equivalent. A signature is equivalent to another, if
// they both have the same scalar value for R and S.
func (sig *Signature) IsEqual(otherSig *Signature) bool {
	return sig.R.Cmp(otherSig.R) == 0 &&
		sig.S.Cmp(otherSig.S) == 0
}

func parseSig(curve elliptic.Curve, sigStr []byte, der bool) (*Signature, error) {
	// Originally this code used encoding/asn1 in order to parse the
	// signature, but a number of problems were found with this approach.
	// Despite the fact that signatures are stored as DER, the difference
	// between go's idea of a bignum (and that they have sign) doesn't agree
	// with the openssl one (where they do not). The above is true as of
	// Go 1.1. In the end it was simpler to rewrite the code to explicitly
	// understand the format which is this:
	// 0x30 <length of whole message> <0x02> <length of R> <R> 0x2
	// <length of S> <S>.

	signature := &Signature{
		curve: curve,
	}

	// minimal message is when both numbers are 1 bytes. adding up to:
	// 0x30 + len + 0x02 + 0x01 + <byte> + 0x2 + 0x01 + <byte>
	if len(sigStr) < 8 {
		return nil, errors.New("malformed signature: too short")
	}
	// 0x30
	index := 0
	if sigStr[index] != 0x30 {
		return nil, errors.New("malformed signature: no header magic")
	}
	index++
	// length of remaining message
	siglen := sigStr[index]
	index++
	if int(siglen+2) > len(sigStr) {
		return nil, errors.New("malformed signature: bad length")
	}
	// trim the slice we're working on so we only look at what matters.
	sigStr = sigStr[:siglen+2]

	// 0x02
	if sigStr[index] != 0x02 {
		return nil,
			errors.New("malformed signature: no 1st int marker")
	}
	index++

	// Length of signature R.
	rLen := int(sigStr[index])
	// must be positive, must be able to fit in another 0x2, <len> <s>
	// hence the -3. We assume that the length must be at least one byte.
	index++
	if rLen <= 0 || rLen > len(sigStr)-index-3 {
		return nil, errors.New("malformed signature: bogus R length")
	}

	// Then R itself.
	rBytes := sigStr[index : index+rLen]
	if der {
		switch err := canonicalPadding(rBytes); err {
		case errNegativeValue:
			return nil, errors.New("signature R is negative")
		case errExcessivelyPaddedValue:
			return nil, errors.New("signature R is excessively padded")
		}
	}
	signature.R = new(big.Int).SetBytes(rBytes)
	index += rLen
	// 0x02. length already checked in previous if.
	if sigStr[index] != 0x02 {
		return nil, errors.New("malformed signature: no 2nd int marker")
	}
	index++

	// Length of signature S.
	sLen := int(sigStr[index])
	index++
	// S should be the rest of the string.
	if sLen <= 0 || sLen > len(sigStr)-index {
		return nil, errors.New("malformed signature: bogus S length")
	}

	// Then S itself.
	sBytes := sigStr[index : index+sLen]
	if der {
		switch err := canonicalPadding(sBytes); err {
		case errNegativeValue:
			return nil, errors.New("signature S is negative")
		case errExcessivelyPaddedValue:
			return nil, errors.New("signature S is excessively padded")
		}
	}
	signature.S = new(big.Int).SetBytes(sBytes)
	index += sLen

	// sanity check length parsing
	if index != len(sigStr) {
		return nil, fmt.Errorf("malformed signature: bad final length %v != %v",
			index, len(sigStr))
	}

	// Verify also checks this, but we can be more sure that we parsed
	// correctly if we verify here too.
	// FWIW the ecdsa spec states that R and S must be | 1, N - 1 |
	// but crypto/ecdsa only checks for Sign != 0. Mirror that.
	if signature.R.Sign() != 1 {
		return nil, errors.New("signature R isn't 1 or more")
	}
	if signature.S.Sign() != 1 {
		return nil, errors.New("signature S isn't 1 or more")
	}
	if signature.R.Cmp(curve.Params().N) >= 0 {
		return nil, errors.New("signature R is >= curve.N")
	}
	if signature.S.Cmp(curve.Params().N) >= 0 {
		return nil, errors.New("signature S is >= curve.N")
	}

	return signature, nil
}

// ParseSignature parses a signature in BER format for the curve type `curve'
// into a Signature type, perfoming some basic sanity checks.  If parsing
// according to the more strict DER format is needed, use ParseDERSignature.
func ParseSignature(curve elliptic.Curve, sigStr []byte) (*Signature, error) {
	return parseSig(curve, sigStr, false)
}

// ParseDERSignature parses a signature in DER format for the curve type
// `curve` into a Signature type.  If parsing according to the less strict
// BER format is needed, use ParseSignature.
func ParseDERSignature(curve elliptic.Curve, sigStr []byte) (*Signature, error) {
	return parseSig(curve, sigStr, true)
}

// canonicalizeInt returns the bytes for the passed big integer adjusted as
// necessary to ensure that a big-endian encoded integer can't possibly be
// misinterpreted as a negative number.  This can happen when the most
// significant bit is set, so it is padded by a leading zero byte in this case.
// Also, the returned bytes will have at least a single byte when the passed
// value is 0.  This is required for DER encoding.
func canonicalizeInt(val *big.Int) []byte {
	b := val.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}
	if b[0]&0x80 != 0 {
		paddedBytes := make([]byte, len(b)+1)
		copy(paddedBytes[1:], b)
		b = paddedBytes
	}
	return b
}

// canonicalPadding checks whether a big-endian encoded integer could
// possibly be misinterpreted as a negative number (even though OpenSSL
// treats all numbers as unsigned), or if there is any unnecessary
// leading zero padding.
func canonicalPadding(b []byte) error {
	switch {
	case b[0]&0x80 == 0x80:
		return errNegativeValue
	case len(b) > 1 && b[0] == 0x00 && b[1]&0x80 != 0x80:
		return errExcessivelyPaddedValue
	default:
		return nil
	}
}

// recoverKeyFromSignature recovers a public key from the signature "sig" on the
// given message hash "msg". Based on the algorithm found in section 5.1.5 of
// SEC 1 Ver 2.0, page 47-48 (53 and 54 in the pdf). This performs the details
// in the inner loop in Step 1. The counter provided is actually the j parameter
// of the loop * 2 - on the first iteration of j we do the R case, else the -R
// case in step 1.6. This counter is used in the Decred compressed signature
// format and thus we match bitcoind's behaviour here.
func recoverKeyFromSignature(sig *Signature, msg []byte, iter int, doChecks bool) (*PublicKey, error) {
	// 1.1 x = (n * i) + r
	curve := sig.curve
	Rx := new(big.Int).Mul(curve.Params().N,
		new(big.Int).SetInt64(int64(iter/2)))
	Rx.Add(Rx, sig.R)
	// 判断有没有超阶
	if Rx.Cmp(curve.Params().P) != -1 {
		return nil, errors.New("calculated Rx is larger than curve P")
	}

	// convert 02<Rx> to point R. (step 1.2 and 1.3). If we are on an odd
	// iteration then 1.6 will be done with -R, so we calculate the other
	// term when uncompressing the point.
	Ry, err := decompressPoint(sig.curve, Rx, iter%2 == 1)
	if err != nil {
		return nil, err
	}

	// 1.4 Check n*R is point at infinity
	if doChecks {
		nRx, nRy := curve.ScalarMult(Rx, Ry, curve.Params().N.Bytes())
		if nRx.Sign() != 0 || nRy.Sign() != 0 {
			return nil, errors.New("n*R does not equal the point at infinity")
		}
	}

	// 1.5 calculate e from message using the same algorithm as ecdsa
	// signature calculation.
	e := hashToInt(sig.curve, msg)

	// Step 1.6.1:
	// We calculate the two terms sR and eG separately multiplied by the
	// inverse of r (from the signature). We then add them to calculate
	// Q = r^-1(sR-eG)
	invr := new(big.Int).ModInverse(sig.R, curve.Params().N)
	if invr == nil {
		invr = big.NewInt(int64(0))
	}
	// first term.
	invrS := new(big.Int).Mul(invr, sig.S)
	invrS.Mod(invrS, curve.Params().N)
	sRx, sRy := curve.ScalarMult(Rx, Ry, invrS.Bytes())

	// second term.
	e.Neg(e)
	e.Mod(e, curve.Params().N)
	e.Mul(e, invr)
	e.Mod(e, curve.Params().N)
	minuseGx, minuseGy := curve.ScalarBaseMult(e.Bytes())

	// step to prevent the jacobian conversion back and forth.
	Qx, Qy := curve.Add(sRx, sRy, minuseGx, minuseGy)

	return &PublicKey{
		Curve: curve,
		X:     Qx,
		Y:     Qy,
	}, nil
}

// SignCompact produces a compact signature of the data in hash with the given
// private key on the given koblitz curve. The isCompressed  parameter should
// be used to detail if the given signature should reference a compressed
// public key or not. If successful the bytes of the compact signature will be
// returned in the format:
// <(byte of 27+public key solution)+4 if compressed >< padded bytes for signature R><padded bytes for signature S>
// where the R and S parameters are padde up to the bitlengh of the curve.
func SignCompact(key *PrivateKey, hash []byte, isCompressedKey bool) ([]byte, error) {
	sig, err := key.Sign(hash)
	if err != nil {
		return nil, err
	}
	privkey := key.ToECDSA()
	curve := privkey.Curve
	// bitcoind checks the bit length of R and S here. The ecdsa signature
	// algorithm returns R and S mod N therefore they will be the bitsize of
	// the curve, and thus correctly sized.
	var H = 1
	for i := 0; i < (H+1)*2; i++ {
		// 通过签名恢复公钥
		pk, err := recoverKeyFromSignature(sig, hash, i, true)
		if err == nil && pk.X.Cmp(key.X) == 0 && pk.Y.Cmp(key.Y) == 0 {
			result := make([]byte, 0, 2*curve.Params().BitSize)
			// 将校验码放在第一位
			//result[0] = 27 + byte(i)
			//if isCompressedKey {
			//	result[0] += 4
			//}
			// Not sure this needs rounding but safer to do so.
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

			// 将校验码放在最后
			verifyCode := byte(i) // add back recid to get 65 bytes sig
			//verifyCode := 27 + byte(i)
			//if isCompressedKey {
			//	verifyCode += 4
			//}
			result = append(result, verifyCode)
			return result, nil
		}
	}

	return nil, errors.New("no valid solution for pubkey found")
}

// RecoverCompact verifies the compact signature "signature" of "hash" for the
// Koblitz curve in "curve". If the signature matches then the recovered public
// key will be returned as well as a boolen if the original key was compressed
// or not, else an error will be returned.
func RecoverCompact(curve elliptic.Curve, signature, hash []byte) (*PublicKey, bool, error) {
	bitlen := (curve.Params().BitSize + 7) / 8
	if len(signature) != 1+bitlen*2 {
		return nil, false, errors.New("invalid compact signature size")
	}

	iteration := int((signature[len(signature)-1]) & ^byte(4))

	// format is <header byte><bitlen R><bitlen S>
	// 校验码在第一位
	//sig := &Signature{
	//	R: new(big.Int).SetBytes(signature[1 : bitlen+1]),
	//	S: new(big.Int).SetBytes(signature[bitlen+1:]),
	//}
	// 校验码在最后一位
	sig := &Signature{
		curve: curve,
		R:     new(big.Int).SetBytes(signature[:bitlen]),
		S:     new(big.Int).SetBytes(signature[bitlen : bitlen*2]),
	}

	// The iteration used here was encoded
	key, err := recoverKeyFromSignature(sig, hash, iteration, false)
	if err != nil {
		return nil, false, err
	}

	return key, ((signature[len(signature)-1]) & 4) == 4, nil
}

// VerifySignature 验证签名
func VerifySignature(pubkey *ecdsa.PublicKey, hash, signature []byte) bool {
	bitlen := (pubkey.Curve.Params().BitSize + 7) / 8
	if len(signature) != bitlen*2 {
		signature = signature[:len(signature)-1]
	}

	// 校验码在最后一位
	sig := &Signature{
		R: new(big.Int).SetBytes(signature[:bitlen]),
		S: new(big.Int).SetBytes(signature[bitlen : bitlen*2]),
	}
	verify := sig.Verify(hash, pubkey)
	return verify
}

// RecoverPubkey 返回未压缩的公钥
func RecoverPubkey(curve elliptic.Curve, hash []byte, signature []byte) ([]byte, error) {
	publicKey, _, err := RecoverCompact(curve, signature, hash)
	if err != nil {
		return nil, err
	}
	return publicKey.SerializeUncompressed(), nil
}
