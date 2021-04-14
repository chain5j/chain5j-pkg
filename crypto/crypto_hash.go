// description: chain5j-core
// 
// @author: xwc1125
// @date: 2020/3/2
package crypto

import (
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/crypto/sha3"
	"github.com/chain5j/chain5j-pkg/types"
)

// keccak256Hash
func RlpHash(x interface{}) (h types.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
