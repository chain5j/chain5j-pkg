// Package secp256k1
//
// @author: xwc1125
package secp256k1

import (
	"crypto/elliptic"

	"github.com/chain5j/chain5j-pkg/crypto/signature/secp256k1/btcecv1"
)

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return btcecv1.S256()
}
