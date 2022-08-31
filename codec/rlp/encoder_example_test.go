// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rlp_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/chain5j/chain5j-pkg/codec/rlp"
)

type MyCoolType struct {
	Name string
	a, b uint
}

// EncodeRLP writes x as RLP list [a, b] that omits the Name field.
func (x *MyCoolType) EncodeRLP(w io.Writer) (err error) {
	return rlp.Encode(w, []uint{x.a, x.b})
}

func ExampleEncoder() {
	var t *MyCoolType // t is nil pointer to MyCoolType
	bytes, _ := rlp.EncodeToBytes(t)
	fmt.Printf("%v → %X\n", t, bytes)

	t = &MyCoolType{Name: "foobar", a: 5, b: 6}
	bytes, _ = rlp.EncodeToBytes(t)
	fmt.Printf("%v → %X\n", t, bytes)

	// Output:
	// <nil> → C0
	// &{foobar 5 6} → C20506
}

type A struct {
	Name string
}
type AA struct {
	Name string
	Age  uint64 `rlp:"optional"`
}
type AAA struct {
	Name string
	Age  uint64 `rlp:"optional"`
	Sex  bool   `rlp:"optional"`
}
type AAAA struct {
	Name string
	Age  uint64 `rlp:"optional"`
	Sex  bool   `rlp:"optional"`
	Data []byte `rlp:"optional"`
}

func TestRlpDecodeA(t *testing.T) {
	a := A{
		Name: "姓名",
	}
	a1Bytes, err := rlp.EncodeToBytes(a)
	if err != nil {
		t.Fatal(err)
	}
	{
		// decode1
		var a A
		err := rlp.DecodeBytes(a1Bytes, &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("decode1", a)
	}
	{
		// decode2
		var a AA
		err := rlp.DecodeBytes(a1Bytes, &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("decode2", a)
	}
	{
		// decode3
		var a AAA
		err := rlp.DecodeBytes(a1Bytes, &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("decode3", a)
	}
	{
		// decode4
		var a AAAA
		err := rlp.DecodeBytes(a1Bytes, &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("decode4", a)
	}
}
func TestRlpDecodeAA(t *testing.T) {
	a := AA{
		Name: "姓名",
		Age:  12,
	}
	a1Bytes, err := rlp.EncodeToBytes(a)
	if err != nil {
		t.Fatal(err)
	}
	{
		// decode2
		var a AA
		err := rlp.DecodeBytes(a1Bytes, &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("decode2", a)
	}
	{
		// decode3
		var a AAA
		err := rlp.DecodeBytes(a1Bytes, &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("decode3", a)
	}
	{
		// decode4
		var a AAAA
		err := rlp.DecodeBytes(a1Bytes, &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("decode4", a)
	}
}
