// Package linkedHashMap
//
// @author: xwc1125
// @date: 2020/12/22
package linkedHashMap

type LinkListNode struct {
	last *LinkListNode
	next *LinkListNode
	val  interface{}
}

func NewLinkListNode(last *LinkListNode, next *LinkListNode, val interface{}) *LinkListNode {
	node := &LinkListNode{
		last: last,
		next: next,
		val:  val,
	}
	return node
}

func (n *LinkListNode) SetLast(node *LinkListNode) {
	n.last = node
}

func (n *LinkListNode) SetNext(node *LinkListNode) {
	n.next = node
}

func (n *LinkListNode) GetLast() *LinkListNode {
	return n.last
}

func (n *LinkListNode) GetNext() *LinkListNode {
	return n.next
}

func (n *LinkListNode) GetVal() interface{} {
	return n.val
}

func (n *LinkListNode) IsHead() bool {
	return n.last == nil
}

func (n *LinkListNode) IsTail() bool {
	return n.next == nil
}

type LinkList struct {
	head   *LinkListNode
	tail   *LinkListNode
	length int
}

func NewLinkList() *LinkList {
	return &LinkList{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

func (ll *LinkList) GetHead() *LinkListNode {
	return ll.head
}

func (ll *LinkList) GetTail() *LinkListNode {
	return ll.tail
}

func (ll *LinkList) AddToHead(val interface{}) *LinkListNode {
	if ll.head == nil && ll.tail == nil {
		return ll.addFirstNode(val)

	}
	node := NewLinkListNode(nil, ll.head, val)
	ll.head.SetLast(node)
	ll.head = node
	ll.length++
	return node
}

func (ll *LinkList) AddToTail(val interface{}) *LinkListNode {
	if ll.head == nil && ll.tail == nil {
		return ll.addFirstNode(val)

	}
	node := NewLinkListNode(ll.tail, nil, val)
	ll.tail.SetNext(node)
	ll.tail = node
	ll.length++
	return node
}

func (ll *LinkList) RemoveNode(node *LinkListNode) {
	defer func() {
		ll.length--
	}()

	/* LinkList中只有1个元素 */
	if node.IsHead() && node.IsTail() {
		ll.head = nil
		ll.tail = nil
		return
	}

	/* 节点是头节点 */
	if node.IsHead() {
		nextNode := node.GetNext()
		ll.head = nextNode
		nextNode.SetLast(nil)
		node.SetNext(nil)
		return
	}

	/* 节点是尾节点 */
	if node.IsTail() {
		lastNode := node.GetLast()
		ll.tail = lastNode
		lastNode.SetNext(nil)
		node.SetLast(nil)
		return
	}

	lastNode := node.GetLast()
	nextNode := node.GetNext()

	lastNode.SetNext(nextNode)
	nextNode.SetLast(lastNode)
	node.SetLast(nil)
	node.SetNext(nil)
}

func (ll *LinkList) GetLength() int {
	return ll.length
}

func (ll *LinkList) addFirstNode(val interface{}) *LinkListNode {
	node := NewLinkListNode(nil, nil, val)
	ll.head = node
	ll.tail = node
	ll.length++
	return node
}
