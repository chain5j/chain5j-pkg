package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/crypto/gmsm"
	"github.com/chain5j/chain5j-pkg/crypto/keccak"
	"github.com/chain5j/chain5j-pkg/crypto/prime256v1"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/types"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

var errInvalidPubkey = errors.New("invalid public key")

// CreateAddress creates an address given the bytes and the nonce
func CreateAddress(b types.Address, nonce *big.Int) types.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return types.BytesToAddress(keccak.Keccak256(data)[12:])
}

// CreateAddress2 creates an address given the address bytes, initial
// contract code and a salt.
func CreateAddress2(b types.Address, salt [32]byte, code []byte) types.Address {
	return types.BytesToAddress(keccak.Keccak256([]byte{0xff}, b.Bytes(), salt[:], keccak.Keccak256(code))[12:])
}

//func CreateAddress3(b types.Address, nonce uint64, payload []byte) types.DomainAddress {
//	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce, payload})
//	return types.BytesToDomainAddress(keccak.Keccak256(data)[12:])
//}

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

// FromECDSA exports a private key into a binary dump.
func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

func MarshalPubkey(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	name := CurveName(pub.Curve)
	switch name {
	case P256:
		return elliptic.Marshal(CurveType(name), pub.X, pub.Y)
	case S256:
		return elliptic.Marshal(CurveType(name), pub.X, pub.Y)
	case SM2P256:
		return gmsm.CompressPubkey(pub)
	}
	return elliptic.Marshal(CurveType(name), pub.X, pub.Y)
}

// UnmarshalPubkey converts bytes to a  public key.
func UnmarshalPubkey(curveName string, pub []byte) (*ecdsa.PublicKey, error) {
	var (
		x, y *big.Int
	)
	switch curveName {
	case P256:
		x, y = elliptic.Unmarshal(CurveType(curveName), pub)
	case S256:
		x, y = elliptic.Unmarshal(CurveType(curveName), pub)
	case SM2P256:
		return gmsm.DecompressPubkey(pub)

	}
	if x == nil {
		return nil, errInvalidPubkey
	}
	return &ecdsa.PublicKey{Curve: CurveType(curveName), X: x, Y: y}, nil
}

// HexToECDSA parses a  private key.
func HexToECDSA(curveName, hexkey string) (*ecdsa.PrivateKey, error) {
	if strings.HasPrefix(hexkey, "0x") {
		hexkey = hexkey[2:]
	}
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToECDSA(curveName, b)
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

func GenerateKey(curveType string) (*ecdsa.PrivateKey, error) {
	switch curveType {
	case P256:
		return prime256v1.GenerateEcdsaPrivateKey()
	case S256:
		break
	case SM2P256:
		return gmsm.GenerateKey()
	}
	return nil, errors.New("curveType is err")
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(big.NewInt(1)) < 0 || s.Cmp(big.NewInt(1)) < 0 {
		return false
	}
	curve := DefaultCryptoType()
	curve256N := curve.Params().N
	curve256halfN := new(big.Int).Div(curve256N, big.NewInt(2))
	if homestead && s.Cmp(curve256halfN) > 0 {
		return false
	}
	// Frontier: allow s to be in full N range
	return r.Cmp(curve256N) < 0 && s.Cmp(curve256N) < 0 && (v == 0 || v == 1)
}

func PubkeyToAddress(p ecdsa.PublicKey) types.Address {
	pubBytes := MarshalPubkey(&p)
	return types.BytesToAddress(keccak.Keccak256(pubBytes[1:])[12:])
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
