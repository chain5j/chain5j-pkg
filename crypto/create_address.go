// Package crypto
//
// @author: xwc1125
// @date: 2019/9/9
package crypto

import (
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/crypto/hash/sha3"
	"github.com/chain5j/chain5j-pkg/types"
	"math/big"
)

var cryptoAlgMap = make(map[string]PrivateKey)

// CreateAddress creates an address given the bytes and the nonce
func CreateAddress(b types.Address, nonce *big.Int) types.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return types.BytesToAddress(sha3.Keccak256(data)[12:])
}

// CreateAddress2 creates an address given the address bytes, initial
// contract code and a salt.
func CreateAddress2(b types.Address, salt [32]byte, code []byte) types.Address {
	return types.BytesToAddress(sha3.Keccak256([]byte{0xff}, b.Bytes(), salt[:], sha3.Keccak256(code))[12:])
}
