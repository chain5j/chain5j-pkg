// Package hashalg
//
// @author: xwc1125
package hashalg

import (
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/crypto/hashalg/sha3"
	"github.com/chain5j/chain5j-pkg/types"
)

// RlpHash keccak256Hash
func RlpHash(x interface{}) (h types.Hash, err error) {
	hw := sha3.NewKeccak256()
	err = rlp.Encode(hw, x)
	if err != nil {
		return types.Hash{}, err
	}
	hw.Sum(h[:0])
	return h, nil
}
