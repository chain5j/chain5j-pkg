// Package reflectutil
// 
// @author: xwc1125
package reflectutil

import (
	"fmt"
	"testing"
)

type AA struct {
	Name  string
	Index int
}

func TestToPointer(t *testing.T) {
	aa := AA{
		Name:  "11",
		Index: 1,
	}
	pointer := ToPointer(aa)
	fmt.Println("===",pointer)

	bb := &AA{
		Name:  "22",
		Index: 2,
	}
	delPointer := DelPointer(bb)
	fmt.Println(delPointer)
}
