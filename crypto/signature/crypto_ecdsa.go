// Package signature
//
// @author: xwc1125
package signature

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/chain5j/chain5j-pkg/crypto"
	"github.com/chain5j/chain5j-pkg/crypto/hashalg/sha3"
	"github.com/chain5j/chain5j-pkg/crypto/signature/gmsm"
	"github.com/chain5j/chain5j-pkg/crypto/signature/prime256v1"
	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/types"
)

var (
	errInvalidPubkey = errors.New("invalid public key")
	errInvalidCurve  = errors.New("unsupported the curve")
)

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

	// The priv.D must < N
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

func PubkeyToAddress(p *ecdsa.PublicKey) types.Address {
	pubBytes, err := MarshalPubkeyWithECDSA(p)
	if err != nil {
		return types.EmptyAddress
	}
	return types.BytesToAddress(sha3.Keccak256(pubBytes[1:])[12:])
}

// GenerateKeyWithECDSA generate the ecdsa key
func GenerateKeyWithECDSA(curveName string) (*ecdsa.PrivateKey, error) {
	if ecdsa, err := GetECDSA(curveName); err != nil {
		return nil, err
	} else {
		return ecdsa.GenerateKey(CurveType(curveName))
	}
}

// SignWithECDSA use ecdsa sign the raw data bytes
func SignWithECDSA(prvKey *ecdsa.PrivateKey, dataBytes []byte) (*SignResult, error) {
	curveName := CurveName(prvKey.Curve)
	if ecdsa, err := GetECDSA(curveName); err != nil {
		return nil, err
	} else {
		hashBytes, _ := ecdsa.HashMsg(curveName, dataBytes)
		signBytes, err := ecdsa.Sign(prvKey, hashBytes)
		if err != nil {
			return nil, err
		}
		pubkeyBytes, err := ecdsa.MarshalPublicKey(&prvKey.PublicKey)
		if err != nil {
			return nil, err
		}
		signResult := &SignResult{
			Name:      curveName,
			PubKey:    pubkeyBytes,
			Signature: signBytes,
		}
		return signResult, nil
	}
}

// VerifyWithECDSA verify the signature by signResult and dataBytes
func VerifyWithECDSA(signResult *SignResult, dataBytes []byte) bool {
	if signResult == nil {
		return false
	}
	if ecdsa, err := GetECDSA(signResult.Name); err != nil {
		return false
	} else {
		pubkey, err := ecdsa.UnmarshalPublicKey(CurveType(signResult.Name), signResult.PubKey)
		if err != nil {
			return false
		}
		hashBytes, _ := ecdsa.HashMsg(signResult.Name, dataBytes)
		return ecdsa.Verify(pubkey, hashBytes, signResult.Signature)
	}
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

func GetECDSA(curveName string) (crypto.ECDSA, error) {
	switch curveName {
	case P256, P384, P521:
		return prime256v1.ECDSA{}, nil
	case S256:
		return secp256k1.Secp251k1{}, nil
	case SM2P256:
		return gmsm.SM2{}, nil
	default:
		return nil, errInvalidCurve
	}
}
