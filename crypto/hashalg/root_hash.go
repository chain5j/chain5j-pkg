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

package hashalg

import (
	"github.com/chain5j/chain5j-pkg/collection/trees/tree"
	"github.com/chain5j/chain5j-pkg/types"
)

// DerivableList derivableList
type DerivableList interface {
	Len() int
	Key(i int) []byte
	Item(i int) []byte
}

// RootHash get list root
func RootHash(list DerivableList) types.Hash {
	trie := new(tree.Trie)
	if list != nil {
		for i := 0; i < list.Len(); i++ {
			trie.Update(list.Key(i), list.Item(i))
		}
	}
	return trie.Hash()
}
