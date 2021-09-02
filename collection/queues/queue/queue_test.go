// Package queue
//
// @author: xwc1125
// @date: 2020/10/20
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

func TestDel(t *testing.T) {
	queue := NewLinkedQueue()
	queue.PushFront("1")
	queue.Remove("1")

	queue.PushBack("2")
	queue.Remove("2")

	queue.PushFront("3")
	queue.PollFront()
	queue.Remove("3")

	queue.PushFront("1")
	queue.PushBack("2")
	queue.PushFront("3")
	queue.PollFront()
	queue.PollBack()
	//queue.Remove("1")
	queue.Remove("3")
	queue.Remove("2")

	queue.Clear()
}
