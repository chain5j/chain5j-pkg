// description: chain5j-core
// 
// @author: xwc1125
// @date: 2020/3/2
package gmsm

import (
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/tjfoc/gmsm/sm3"
)

// GM3摘要算法
func Gm3Hash(data []byte) (h types.Hash) {
	bytes := gm3Hash(data)
	return types.BytesToHash(bytes)
}

// SM3密码杂凑算法
func gm3Hash(data []byte) []byte {
	h := sm3.New()
	h.Write(data)
	sum := h.Sum(nil)
	return sum
}