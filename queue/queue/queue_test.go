// description: chain5j-pkg
// 
// @author: xwc1125
// @date: 2019/11/25
package queue

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

type item struct {
	nonce int64
	value interface{}
}

func TestNew(t *testing.T) {
	linkedQueue := NewLinkedQueue()

	for i := int64(0); i < 5; i++ {
		linkedQueue.PushFront(Element(item{
			i,
			i,
		}))
	}
	for i := 0; i < 2; i++ {
		spew.Dump(linkedQueue.PeekFront())
		spew.Dump(linkedQueue.PeekBack())
		linkedQueue.PollFront()
	}
	spew.Dump(linkedQueue.PeekFront())
	spew.Dump(linkedQueue.PeekBack())
	linkedQueue.PollBack()
	for i := 0; i < 2; i++ {
		spew.Dump(linkedQueue.PeekFront())
		spew.Dump(linkedQueue.PeekBack())
		linkedQueue.PollFront()
	}
	spew.Dump(linkedQueue.PeekFront())
	spew.Dump(linkedQueue.PeekBack())
	linkedQueue.PollBack()
	spew.Dump(linkedQueue.PeekFront())
	spew.Dump(linkedQueue.PeekBack())
	linkedQueue.PushBack(Element(item{
		8,
		8,
	}))
	spew.Dump(linkedQueue.PeekFront())
	spew.Dump(linkedQueue.PeekBack())
	linkedQueue.PushBack(Element(item{
		9,
		9,
	}))
	spew.Dump(linkedQueue.PeekFront())
	spew.Dump(linkedQueue.PeekBack())
	linkedQueue.PollBack()
	spew.Dump(linkedQueue.PeekFront())
	spew.Dump(linkedQueue.PeekBack())
	linkedQueue.PollBack()
	spew.Dump(linkedQueue.PeekFront())
	spew.Dump(linkedQueue.PeekBack())

	linkedQueue.Clear()

	spew.Dump(linkedQueue.PeekFront())
}

