// description: chain5j-pkg
// 
// @author: xwc1125
// @date: 2019/11/25
package queue

import (
	"fmt"
	"log"
	"sync"
)

type Element interface{}

type Queue interface {
	PushFront(e Element)   // 向队头添加元素
	PushBack(e Element)    // 向队尾添加元素
	PeekFront() Element    // 查看头部的元素
	PeekBack() Element     // 查看尾部的元素
	PollFront() Element    // 移除头部的元素
	PollBack() Element     // 移除尾部的元素
	Remove(e Element) bool // 删除一个值
	Exist(e Element) bool  // 是否存在
	Size() int             // 获取队列的元素个数
	IsEmpty() bool         // 判断队列是否是空
	Clear() bool           // 清空队列
	NewIterator() *Iterator
}

type node struct {
	value Element // 当前节点的值
	prev  *node   // 前一个节点
	next  *node   // 下一个节点
}

type LinkedQueue struct {
	m    sync.Mutex
	head *node // 头节点
	tail *node // 尾节点
	size int   // 大小
}

func NewLinkedQueue() *LinkedQueue {
	return &LinkedQueue{
		size: 0,
	}
}

func (queue *LinkedQueue) PushFront(e Element) {
	queue.m.Lock()
	defer queue.m.Unlock()

	newNode := &node{e, nil, queue.head}
	if queue.head == nil {
		queue.head = newNode
		queue.tail = newNode
	} else {
		queue.head.prev = newNode
		queue.head = newNode
	}
	queue.size++
	newNode = nil
}

func (queue *LinkedQueue) PushBack(e Element) {
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

func (queue *LinkedQueue) PeekFront() Element {
	if queue.head == nil {
		return nil
	}
	return queue.head.value
}

func (queue *LinkedQueue) PeekBack() Element {
	if queue.tail == nil {
		return nil
	}
	return queue.tail.value
}

// 移除队列中最前面的元素
func (queue *LinkedQueue) PollFront() Element {
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

func (queue *LinkedQueue) PollBack() Element {
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

func (queue *LinkedQueue) Remove(e Element) bool {
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
	if cNode == nil {
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
			log.Println("LinkedQueue exist", "err", err)
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
