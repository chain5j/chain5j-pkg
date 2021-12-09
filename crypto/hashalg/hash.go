// Package hashalg
//
// @author: xwc1125
package hashalg

import (
	"crypto/sha256"
)

// Sha256 sha256
func Sha256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
