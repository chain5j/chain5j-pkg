package prque

import "testing"

func TestPrque_PopBottom(t *testing.T) {
	prque := New(nil)
	for i := int64(0); i < 100; i++ {
		prque.Push(i, i)
	}
	prque.Push("2", 2)
	prque.Push("3", 3)
	prque.Push("1", 1)
	// priority:on the top
	// 3
	// 2
	// 1

	peek, priority := prque.PeekTop()
	println("PeekTop", peek, priority)
	peek, priority = prque.PeekBottom()
	println("PeekBottom", peek, priority)
	println("===========================")

	peek, priority = prque.PopBottom()
	println("PopBottom", peek, priority)
	peek, priority = prque.PeekTop()
	println("PeekTop", peek, priority)
	peek, priority = prque.PeekBottom()
	println("PeekBottom", peek, priority)
	println("===========================")

	peek, priority = prque.PopBottom()
	println("PopBottom", peek, priority)
	peek, priority = prque.PeekTop()
	println("PeekTop", peek, priority)
	peek, priority = prque.PeekBottom()
	println("PeekBottom", peek, priority)
	println("===========================")

}

func TestPrque_Pop(t *testing.T) {
	prque := New(nil)
	prque.Push("2", 2)
	prque.Push("3", 3)
	prque.Push("1", 1)
	// priority:on the top
	// 3
	// 2
	// 1

	peek, priority := prque.PeekTop()
	println("PeekTop", peek, priority)
	peek, priority = prque.PeekBottom()
	println("PeekBottom", peek, priority)
	println("===========================")

	peek, priority = prque.Pop()
	println("Pop", peek, priority)
	peek, priority = prque.PeekTop()
	println("PeekTop", peek, priority)
	peek, priority = prque.PeekBottom()
	println("PeekBottom", peek, priority)
	println("===========================")

	peek, priority = prque.Pop()
	println("Pop", peek, priority)
	peek, priority = prque.PeekTop()
	println("PeekTop", peek, priority)
	peek, priority = prque.PeekBottom()
	println("PeekBottom", peek, priority)
	println("===========================")

	peek, priority = prque.Pop()
	println("Pop", peek, priority)
	peek, priority = prque.PeekTop()
	println("PeekTop", peek, priority)
	peek, priority = prque.PeekBottom()
	println("PeekBottom", peek, priority)
	println("===========================")

}
