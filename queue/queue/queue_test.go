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
		linkedQueue.Push(Element(item{
			i,
			i,
		}))
	}
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.Poll()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.Poll()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.PollBottom()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.Poll()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.Poll()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.PollBottom()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.Push(Element(item{
		8,
		8,
	}))
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.Push(Element(item{
		9,
		9,
	}))
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.PollBottom()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())
	linkedQueue.PollBottom()
	spew.Dump(linkedQueue.Peek())
	spew.Dump(linkedQueue.PeekBottom())


	//linkedQueue.Push(Element("1"))
	//linkedQueue.Push(Element("2"))
	//linkedQueue.Push(Element("3"))
	//spew.Dump(linkedQueue.Peek())
	//element := linkedQueue.Poll()
	//spew.Dump(element)
	//spew.Dump(linkedQueue.Size())
	//spew.Dump(linkedQueue.IsEmpty())
	//linkedQueue.Push(Element("5"))
	//
	//spew.Dump(linkedQueue.Peek())
	linkedQueue.Clear()

	spew.Dump(linkedQueue.Peek())
}
