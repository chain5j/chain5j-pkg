// Package secp256k1
//
// @author: xwc1125
package secp256k1

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1/btcecv1"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

// RecoverPubkey returns the uncompressed public key that created the given signature.
func RecoverPubkey(hash, sig []byte) ([]byte, error) {
	pub, err := SigToPub(hash, sig)
	if err != nil {
		return nil, err
	}
	p := (*btcecv1.PublicKey)(pub)
	bytes := p.SerializeUncompressed()
	return bytes, err
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	// Convert to btcec input format with 'recovery id' v at the beginning.
	btcsig := make([]byte, 65)
	btcsig[0] = sig[64] + 27
	copy(btcsig[1:], sig)

	pub, _, err := btcecv1.RecoverCompact(btcecv1.S256(), btcsig, hash)
	return pub.ToECDSA(), err
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given hash cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(prv *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != btcecv1.S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, err := btcecv1.SignCompact(btcecv1.S256(), (*btcecv1.PrivateKey)(prv), hash, false)
	if err != nil {
		return nil, err
	}
	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

func Verify(pub *ecdsa.PublicKey, hash []byte, signature []byte) bool {
	sig := &btcecv1.Signature{R: new(big.Int).SetBytes(signature[:32]), S: new(big.Int).SetBytes(signature[32:64])}
	key := (*btcecv1.PublicKey)(pub)
	return sig.Verify(hash, key)
}

// VerifySignature checks that the given public key created signature over hash.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature(pubkey, hash, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}
	sig := &btcecv1.Signature{R: new(big.Int).SetBytes(signature[:32]), S: new(big.Int).SetBytes(signature[32:])}
	key, err := btcecv1.ParsePubKey(pubkey, btcecv1.S256())
	if err != nil {
		return false
	}
	return sig.Verify(hash, key)
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	if len(pubkey) != 33 {
		return nil, errors.New("invalid compressed public key length")
	}
	key, err := btcecv1.ParsePubKey(pubkey, btcecv1.S256())
	if err != nil {
		return nil, err
	}
	return key.ToECDSA(), nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	return (*btcecv1.PublicKey)(pubkey).SerializeCompressed()
}
