// This is a duplicated and slightly modified version of "gopkg.in/karalabe/cookiejar.v2/collections/prque".

package prque

// The size of a block of data
//const blockSize = 4096
const blockSize = 5

// A prioritized item in the sorted stack.
//
// Note: priorities can "wrap around" the int64 range, a comes before b if (a.priority - b.priority) > 0.
// The difference between the lowest and highest priorities in the queue at any point should be less than 2^63.
type item struct {
	value    interface{}
	priority int64
}

// setIndexCallback is called when the element is moved to a new index.
// Providing setIndexCallback is optional, it is needed only if the application needs
// to delete elements other than the top one.
type setIndexCallback func(a interface{}, i int)

// Internal sortable stack data structure. Implements the Push and Pop ops for
// the stack (heap) functionality and the Len, Less and Swap methods for the
// sortability requirements of the heaps.
type sstack struct {
	setIndex setIndexCallback
	size     int
	capacity int
	offset   int

	blocks [][]*item
	active []*item
}

// Creates a new, empty stack.
func newSstack(setIndex setIndexCallback) *sstack {
	result := new(sstack)
	result.setIndex = setIndex
	result.active = make([]*item, blockSize)
	result.blocks = [][]*item{result.active}
	result.capacity = blockSize
	return result
}

// Pushes a value onto the stack, expanding it if necessary. Required by
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
	if s.setIndex != nil {
		s.setIndex(data.(*item).value, s.size)
	}
	s.active[s.offset] = data.(*item)
	s.offset++
	s.size++
}

// Pops a value off the stack and returns it. Currently no shrinking is done.
// Required by heap.Interface.
func (s *sstack) Pop() (res interface{}) {
	s.size--
	s.offset--
	if s.offset < 0 {
		s.offset = blockSize - 1
		s.active = s.blocks[s.size/blockSize]
	}
	res, s.active[s.offset] = s.active[s.offset], nil
	if s.setIndex != nil {
		s.setIndex(res.(*item).value, -1)
	}
	return
}

func (s *sstack) PopBottom() (res interface{}) {
	s.size--
	s.offset--
	if s.offset < 0 {
		s.offset = blockSize - 1
		s.active = s.blocks[s.size/blockSize]
	}
	res, s.active[s.offset] = s.active[s.offset], nil
	return res
}

// Returns the length of the stack. Required by sort.Interface.
func (s *sstack) Len() int {
	return s.size
}

// Compares the priority of two elements of the stack (higher is first).
// Required by sort.Interface.
func (s *sstack) Less(i, j int) bool {
	return (s.blocks[i/blockSize][i%blockSize].priority - s.blocks[j/blockSize][j%blockSize].priority) > 0
}

// Swaps two elements in the stack. Required by sort.Interface.
func (s *sstack) Swap(i, j int) {
	ib, io, jb, jo := i/blockSize, i%blockSize, j/blockSize, j%blockSize
	a, b := s.blocks[jb][jo], s.blocks[ib][io]
	if s.setIndex != nil {
		s.setIndex(a.value, i)
		s.setIndex(b.value, j)
	}
	s.blocks[ib][io], s.blocks[jb][jo] = a, b
}

// Resets the stack, effectively clearing its contents.
func (s *sstack) Reset() {
	*s = *newSstack(s.setIndex)
}

func (s *sstack) PeekTop() (data interface{}, priority int64) {
	offset := s.offset
	active := s.active
	if len(s.blocks) > 1 {
		active = s.blocks[0]
	}
	if offset < 0 {
		offset = blockSize - 1
		active = s.blocks[s.size/blockSize]
	}
	item := active[0]
	if item == nil {
		return nil, 0
	}
	return item.value, item.priority
}

func (s *sstack) PeekBottom() (data interface{}, priority int64) {
	if s.offset == 0 {
		return nil, 0
	}
	offset := s.offset
	active := s.active
	if s.offset < 0 {
		offset = blockSize - 1
		active = s.blocks[s.size/blockSize]
	}
	item := active[offset-1]
	if item == nil {
		return nil, 0
	}
	return item.value, item.priority
}
