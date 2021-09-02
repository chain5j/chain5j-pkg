// Package sha3
//
// @author: xwc1125
// @date: 2020/10/11
package sha3

import (
	"github.com/chain5j/chain5j-pkg/types"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h types.Hash) {
	d := NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

// Keccak512 calculates and returns the Keccak512 hash of the input data.
func Keccak512(data ...[]byte) []byte {
	d := NewKeccak512()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
