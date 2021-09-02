// Package base36
//
// @author: xwc1125
// @date: 2021/4/14
package base36

import (
	"fmt"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	base36Str := EncodeBytes([]byte("0x771a45c3e50550878489f9361b7ba94a96c0ced7"))
	fmt.Println(base36Str)

	bytes := DecodeToBytes(true, base36Str)
	fmt.Println("upper", string(bytes))

	bytes2 := DecodeToBytes(false, strings.ToLower(base36Str))
	fmt.Println("lower", string(bytes2))
}
