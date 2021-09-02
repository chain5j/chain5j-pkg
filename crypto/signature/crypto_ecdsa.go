// Package signature
//
// @author: xwc1125
// @date: 2021/8/10
package signature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/chain5j/chain5j-pkg/crypto"
	"github.com/chain5j/chain5j-pkg/crypto/hash/sha3"
	"github.com/chain5j/chain5j-pkg/crypto/signature/gmsm"
	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/tjfoc/gmsm/sm3"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

var errInvalidPubkey = errors.New("invalid public key")

// ToECDSA creates a private key with the given D value.
func ToECDSA(curveName string, d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(curveName, d, true)
}

// ToECDSAUnsafe blindly converts a binary blob to a private key. It should almost
// never be used unless you are sure the input is valid and want to avoid hitting
// errors due to bad origin encoding (0 prefixes cut off).
func ToECDSAUnsafe(curveName string, d []byte) *ecdsa.PrivateKey {
	priv, _ := toECDSA(curveName, d, false)
	return priv
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA(curveName string, d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)

	priv.PublicKey.Curve = CurveType(curveName)
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N,secp256k1N
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

// ToHexWithECDSA convert ecdsa privateKey to hex
func ToHexWithECDSA(prvKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(FromECDSA(prvKey))
}

// HexToECDSA parses a  private key.
func HexToECDSA(curveName, hexKey string) (*ecdsa.PrivateKey, error) {
	if strings.HasPrefix(hexKey, "0x") {
		hexKey = hexKey[2:]
	}
	b, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToECDSA(curveName, b)
}

// FromECDSA exports a private key into a binary dump.
func FromECDSA(prvKey *ecdsa.PrivateKey) []byte {
	if prvKey == nil {
		return nil
	}
	return math.PaddedBigBytes(prvKey.D, prvKey.Params().BitSize/8)
}

// MarshalPubkeyWithECDSA marshal ecdsa publicKey to bytes
func MarshalPubkeyWithECDSA(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	name := CurveName(pub.Curve)
	switch name {
	case P256, P384, P521:
		return elliptic.Marshal(CurveType(name), pub.X, pub.Y)
	case S256:
		return elliptic.Marshal(CurveType(name), pub.X, pub.Y)
	case SM2P256:
		return gmsm.MarshalPublicKeyFromECDSA(pub)
	}
	return elliptic.Marshal(CurveType(name), pub.X, pub.Y)
}

// UnmarshalPubkeyWithECDSA converts bytes to a  public key.
func UnmarshalPubkeyWithECDSA(curveName string, pub []byte) (*ecdsa.PublicKey, error) {
	if len(pub) == 0 {
		return nil, errors.New("pubBytes is empty")
	}
	var (
		x, y *big.Int
	)
	switch curveName {
	case P256, P384, P521:
		curve := CurveType(curveName)
		x, y = elliptic.Unmarshal(curve, pub)
		if x == nil {
			return nil, errInvalidPubkey
		}
		return &ecdsa.PublicKey{X: x, Y: y, Curve: curve}, nil
	case S256:
		curve := CurveType(curveName)
		x, y = elliptic.Unmarshal(curve, pub)
		if x == nil {
			return nil, errInvalidPubkey
		}
		return &ecdsa.PublicKey{X: x, Y: y, Curve: curve}, nil
	case SM2P256:
		return gmsm.UnmarshalPublicKeyToECDSA(pub)

	}
	return nil, errors.New("unsupported the curveName")
}

// LoadECDSA loads a  private key from the given file.
func LoadECDSA(curveName string, file string) (*ecdsa.PrivateKey, error) {
	buf := make([]byte, 64)
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err := io.ReadFull(fd, buf); err != nil {
		return nil, err
	}

	key, err := hex.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	return ToECDSA(curveName, key)
}

// SaveECDSA saves a  private key to the given file with
// restrictive permissions. The key data is saved hex-encoded.
func SaveECDSA(file string, key *ecdsa.PrivateKey) error {
	k := hex.EncodeToString(FromECDSA(key))
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// GenerateKeyWithECDSA generate the ecdsa key
func GenerateKeyWithECDSA(curveName string) (*ecdsa.PrivateKey, error) {
	switch curveName {
	case P256, P384, P521:
		return ecdsa.GenerateKey(CurveType(curveName), rand.Reader)
	case S256:
		return secp256k1.GenerateKey()
	case SM2P256:
		return gmsm.GenerateKey()
	}
	return nil, errors.New("unsupported curve name")
}

func MarshalPrvkeyWithECDSA(prv *ecdsa.PrivateKey) ([]byte, error) {
	if prv == nil {
		return nil, errors.New("ecdsa private key is empty")
	}
	name := CurveName(prv.Curve)
	switch name {
	case P256, P384, P521, S256:
		return x509.MarshalECPrivateKey(prv)
	case SM2P256:
		return gmsm.MarshalPrivateKey(prv)
	}
	return nil, errors.New("unsupported the CurveName")
}

func PubkeyToAddress(p ecdsa.PublicKey) types.Address {
	pubBytes := MarshalPubkeyWithECDSA(&p)
	return types.BytesToAddress(sha3.Keccak256(pubBytes[1:])[12:])
}

// SignWithECDSA use ecdsa sign the raw data bytes
func SignWithECDSA(prvKey *ecdsa.PrivateKey, dataBytes []byte) (*SignResult, error) {
	curveName := CurveName(prvKey.Curve)
	switch curveName {
	case P256, P384, P521:
		hash := crypto.HashMsg(curveName, dataBytes)
		signBytes, err := prvKey.Sign(rand.Reader, hash, crypto.HashType(curveName))
		if err != nil {
			return nil, err
		}
		signResult := &SignResult{
			Name:      curveName,
			PubKey:    MarshalPubkeyWithECDSA(&prvKey.PublicKey),
			Signature: signBytes,
		}
		return signResult, nil
	case S256:
	case SM2P256:
		hashBytes := sm3.Sm3Sum(dataBytes)
		signBytes, err := gmsm.Sign(hashBytes[:], prvKey)
		if err != nil {
			return nil, err
		}
		signResult := &SignResult{
			Name:      curveName,
			PubKey:    gmsm.MarshalPublicKeyFromECDSA(&prvKey.PublicKey),
			Signature: signBytes,
		}
		return signResult, nil
	}
	return nil, errors.New("unsupported the curve")
}

// VerifyWithECDSA verify the signature by signResult and dataBytes
func VerifyWithECDSA(signResult *SignResult, dataBytes []byte) bool {
	if signResult == nil {
		return false
	}
	switch signResult.Name {
	case P256, P384, P521:
		pubkey, err := UnmarshalPubkeyWithECDSA(signResult.Name, signResult.PubKey)
		if err != nil {
			return false
		}
		hashBytes := crypto.HashMsg(signResult.Name, dataBytes)
		return ecdsa.VerifyASN1(pubkey, hashBytes, signResult.Signature)
	case S256:
	case SM2P256:
		unmarshalPubkey, err := gmsm.UnmarshalPublicKeyToECDSA(signResult.PubKey)
		if err != nil {
			return false
		}
		hashBytes := crypto.HashMsg(signResult.Name, dataBytes)
		return gmsm.VerifyECDSA(unmarshalPubkey, hashBytes[:], signResult.Signature)
	}
	return false
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(big.NewInt(1)) < 0 || s.Cmp(big.NewInt(1)) < 0 {
		return false
	}
	// todo
	curve := CurveType(S256)
	curve256N := curve.Params().N
	curve256halfN := new(big.Int).Div(curve256N, big.NewInt(2))
	if homestead && s.Cmp(curve256halfN) > 0 {
		return false
	}
	// Frontier: allow s to be in full N range
	return r.Cmp(curve256N) < 0 && s.Cmp(curve256N) < 0 && (v == 0 || v == 1)
}
