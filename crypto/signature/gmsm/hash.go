// Package gmsm
//
// @author: xwc1125
package gmsm

import (
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/tjfoc/gmsm/sm3"
)

// Gm3Hash GM3摘要算法
func Gm3Hash(data []byte) (h types.Hash) {
	bytes := Gm3HashBytes(data)
	return types.BytesToHash(bytes)
}

// Gm3HashBytes SM3密码杂凑算法
func Gm3HashBytes(data []byte) []byte {
	h := sm3.New()
	h.Write(data)
	sum := h.Sum(nil)
	return sum
}
