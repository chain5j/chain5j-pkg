// description: chain5j-pkg
// 
// @author: xwc1125
// @date: 2019/11/25
package queue

import (
	"fmt"
	log "github.com/chain5j/log15"
	"sync"
)

type Element interface{}

type Queue interface {
	Push(e Element)        // 向队尾添加元素
	Peek() Element         // 查看头部的元素
	PeekBottom() Element   // 查看尾部的元素
	Poll() Element         // 移除头部的元素
	PollBottom() Element   // 移除尾部的元素
	Delete(e Element) bool // 删除一个值
	Exist(e Element) bool  // 是否存在
	Size() int             // 获取队列的元素个数
	IsEmpty() bool         // 判断队列是否是空
	Clear() bool           // 清空队列
}

type node struct {
	value Element
	prev  *node
	next  *node
}

type LinkedQueue struct {
	head *node
	tail *node
	size int
	m    sync.Mutex
}

func NewLinkedQueue() *LinkedQueue {
	return &LinkedQueue{
		size: 0,
	}
}

func (queue *LinkedQueue) Push(e Element) {
	queue.m.Lock()
	defer queue.m.Unlock()

	newNode := &node{e, queue.tail, nil}
	if queue.tail == nil {
		queue.head = newNode
		queue.tail = newNode
	} else {
		queue.tail.next = newNode
		queue.tail = newNode
	}
	queue.size++
	newNode = nil
}

func (queue *LinkedQueue) Peek() Element {
	if queue.head == nil {
		return nil
	}
	return queue.head.value
}

func (queue *LinkedQueue) PeekBottom() Element {
	if queue.tail == nil {
		return nil
	}
	return queue.tail.value
}

// 移除队列中最前面的元素
func (queue *LinkedQueue) Poll() Element {
	queue.m.Lock()
	defer queue.m.Unlock()
	if queue.IsEmpty() {
		return nil
	}
	if queue.head == nil {
		fmt.Println("Poll Empty queue.")
		return nil
	}
	queue.size--

	firstNode := queue.head
	queue.head = firstNode.next
	if queue.head != nil {
		queue.head.prev = nil
	} else {
		queue.tail = nil
	}

	return firstNode.value
}

func (queue *LinkedQueue) PollBottom() Element {
	queue.m.Lock()
	defer queue.m.Unlock()
	if queue.IsEmpty() {
		return nil
	}
	if queue.tail == nil {
		fmt.Println("PollBottom Empty queue.")
		return nil
	}
	queue.size--

	latestNode := queue.tail
	queue.tail = latestNode.prev
	if queue.tail != nil {
		queue.tail.prev = nil
		queue.tail.next = nil
	} else {
		queue.head = nil
	}

	return latestNode.value
}

func (queue *LinkedQueue) Delete(e Element) bool {
	queue.m.Lock()
	defer queue.m.Unlock()
	if queue.IsEmpty() {
		return true
	}
	if queue.head == nil {
		return true
	}
	firstNode := queue.head
	queue.del(firstNode, e)
	return true
}

func (queue *LinkedQueue) del(cNode *node, e Element) *node {
	if queue.IsEmpty() {
		return nil
	}
	prev := cNode.prev
	next2 := cNode.next
	if cNode.value == e {
		// 查找元素
		prev.next = next2
		next2.prev = prev
		queue.size--
		return cNode
	}
	return queue.del(next2, e)
}

func (queue *LinkedQueue) Exist(e Element) bool {
	queue.m.Lock()
	defer queue.m.Unlock()
	if queue.IsEmpty() {
		return false
	}
	if queue.head == nil {
		return false
	}
	firstNode := queue.head
	return queue.exist(firstNode, e)
}

func (queue *LinkedQueue) exist(cNode *node, e Element) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Error("LinkedQueue exist", "err", err)
		}
	}()
	if queue.IsEmpty() {
		return false
	}
	next2 := cNode.next
	if cNode.value == e {
		// 查找元素
		return true
	}
	return queue.exist(next2, e)
}

func (queue *LinkedQueue) Size() int {
	return queue.size
}

func (queue *LinkedQueue) IsEmpty() bool {
	if queue.size == 0 {
		return true
	}
	return false
}

func (queue *LinkedQueue) Clear() bool {
	if queue.IsEmpty() {
		return false
	}
	queue.m.Lock()
	defer queue.m.Unlock()

	queue.remove()
	return true
}

func (queue *LinkedQueue) remove() {
	if !queue.IsEmpty() {
		firstNode := queue.head
		queue.head = firstNode.next
		firstNode.next = nil
		firstNode.value = nil
		queue.size--

		queue.remove()
	}
}

// NewIterator creates a new iterator for the cache.
func (queue *LinkedQueue) NewIterator() *Iterator {
	return &Iterator{
		current: queue.head,
	}
}
