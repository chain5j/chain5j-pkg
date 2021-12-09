// CookieJar - A contestant's algorithm toolbox
// Copyright (c) 2013 Peter Szilagyi. All rights reserved.
//
// CookieJar is dual licensed: use of this source code is governed by a BSD
// license that can be found in the LICENSE file. Alternatively, the CookieJar
// toolbox may be used in accordance with the terms and conditions contained
// in a signed written agreement between you and the author(s).

package prque

// The size of a block of data
const blockSize = 4096

// A prioritized item in the sorted stack.
type item struct {
	value    interface{}
	priority float32
}

// Internal sortable stack data structure. Implements the Push and Pop ops for
// the stack (heap) functionality and the Len, Less and Swap methods for the
// sortability requirements of the heaps.
type sstack struct {
	size     int
	capacity int
	offset   int

	blocks [][]*item
	active []*item
}

// Creates a new, empty stack.
func newSstack() *sstack {
	result := new(sstack)
	result.active = make([]*item, blockSize)
	result.blocks = [][]*item{result.active}
	result.capacity = blockSize
	return result
}

// Push a value onto the stack, expanding it if necessary. Required by
// heap.Interface.
func (s *sstack) Push(data interface{}) {
	if s.size == s.capacity {
		s.active = make([]*item, blockSize)
		s.blocks = append(s.blocks, s.active)
		s.capacity += blockSize
		s.offset = 0
	} else if s.offset == blockSize {
		s.active = s.blocks[s.size/blockSize]
		s.offset = 0
	}
	s.active[s.offset] = data.(*item)
	s.offset++
	s.size++
}

// Pop a value off the stack and returns it. Currently no shrinking is done.
// Required by heap.Interface.
func (s *sstack) Pop() (res interface{}) {
	s.size--
	s.offset--
	if s.offset < 0 {
		s.offset = blockSize - 1
		s.active = s.blocks[s.size/blockSize]
	}
	res, s.active[s.offset] = s.active[s.offset], nil
	return
}
func (s *sstack) PeekFront() (res interface{}) {
	index := s.offset - 1
	if index < 0 {
		return nil
	}
	res = s.active[index]
	return
}

func (s *sstack) Len() int {
	return s.size
}
func (s *sstack) Less(i, j int) bool {
	return s.blocks[i/blockSize][i%blockSize].priority > s.blocks[j/blockSize][j%blockSize].priority
}
func (s *sstack) Swap(i, j int) {
	ib, io, jb, jo := i/blockSize, i%blockSize, j/blockSize, j%blockSize
	s.blocks[ib][io], s.blocks[jb][jo] = s.blocks[jb][jo], s.blocks[ib][io]
}
func (s *sstack) Reset() {
	*s = *newSstack()
}
