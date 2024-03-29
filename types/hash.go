// Package types
//
// @author: xwc1125
package types

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"

	"github.com/chain5j/chain5j-pkg/crypto/hashalg/sha3"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

// HashLength hash长度
const HashLength = 32

var (
	hashT         = reflect.TypeOf(Hash{})
	EmptyHash     = Hash{}
	EmptyCode     = sha3.Keccak256(nil)
	EmptyCodeHash = BytesToHash(EmptyCode)
	EmptyRootHash = HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	// EmptyRootHash, _ = hashalg.RootHash(codec.Codecor(), nil)
)

type Hash [HashLength]byte

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func StringToHash(s string) Hash {
	return BytesToHash([]byte(s))
}

func BigToHash(b *big.Int) Hash {
	return BytesToHash(b.Bytes())
}

func HexToHash(s string) Hash {
	return BytesToHash(hexutil.MustDecode(s))
}

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

func (h Hash) HexNoPrefix() string {
	return hexutil.Encode(h[:])[2:]
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	// return fmt.Sprintf("%x…%x", h[:3], h[29:])
	return h.Hex()
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (h Hash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h[:])
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

func IsEmptyHash(h Hash) bool {
	return h == Hash{}
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != HashLength {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), HashLength)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

func (h Hash) Nil() bool {
	return h == Hash{}
}

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}
