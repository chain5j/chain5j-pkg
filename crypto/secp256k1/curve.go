// description: chain5j-core
// 
// @author: xwc1125
// @date: 2020/8/8
package secp256k1

import (
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcec"
)

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return btcec.S256()
}
