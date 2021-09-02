// Package queue
//
// @author: xwc1125
// @date: 2020/3/20
package queue

type Iterator struct {
	current *node
}

func (i *Iterator) Value() Element {
	if i.current == nil {
		return nil
	}
	return i.current.value
}

func (i *Iterator) Next() *Iterator {
	if i.current == nil {
		return nil
	}
	i.current = i.current.next
	return i
}
