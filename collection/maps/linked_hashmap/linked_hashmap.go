// Package linkedHashMap
//
// @author: xwc1125
package linkedHashMap

import (
	"container/list"
	"sync"

	"github.com/chain5j/logger"
)

type Node struct {
	node *list.Element
	val  interface{}
}

type LinkedHashMap struct {
	linklistLock *sync.RWMutex
	linklist     *list.List
	hashmap      *sync.Map
}

func NewLinkedHashMap() *LinkedHashMap {
	return &LinkedHashMap{
		linklist:     list.New(),
		linklistLock: new(sync.RWMutex),
		hashmap:      new(sync.Map),
	}
}

func (m *LinkedHashMap) Lock() {
	m.linklistLock.Lock()
}

func (m *LinkedHashMap) Unlock() {
	m.linklistLock.Unlock()
}

func (m *LinkedHashMap) RLock() {
	m.linklistLock.RLock()
}

func (m *LinkedHashMap) RUnlock() {
	m.linklistLock.RUnlock()
}

// Add 添加
func (m *LinkedHashMap) Add(key interface{}, val interface{}) bool {
	_, isExists := m.hashmap.Load(key)
	if isExists {
		return false
	}

	m.Lock()
	linkListNode := m.linklist.PushBack(key)
	m.Unlock()
	m.hashmap.Store(key, &Node{
		node: linkListNode,
		val:  val,
	})

	return true
}
func (m *LinkedHashMap) AddFront(key interface{}, val interface{}) bool {
	_, isExists := m.hashmap.Load(key)
	if isExists {
		return false
	}

	m.Lock()
	linkListNode := m.linklist.PushFront(key)
	m.Unlock()
	m.hashmap.Store(key, &Node{
		node: linkListNode,
		val:  val,
	})

	return true
}

// Get 获取数据
func (m *LinkedHashMap) Get(key interface{}) (interface{}, bool) {
	node, isExists := m.hashmap.Load(key)
	if !isExists {
		return nil, false
	}

	return &node.(*Node).val, true
}

// Exist 判断是否存在
func (m *LinkedHashMap) Exist(key interface{}) bool {
	_, isExists := m.hashmap.Load(key)
	return isExists
}

// Remove 删除
func (m *LinkedHashMap) Remove(key interface{}) {
	node, isExists := m.hashmap.Load(key)
	if !isExists {
		return
	}

	m.linklistLock.Lock()
	m.hashmap.Delete(key)
	m.linklistLock.Unlock()
	m.linklist.Remove(node.(*Node).node)
	return
}

// Len 获取len
func (m *LinkedHashMap) Len() int {
	return m.linklist.Len()
}

// =========批量处理=========

// BatchAdd 批量添加
func (m *LinkedHashMap) BatchAdd(kvs ...KV) {
	if kvs == nil || len(kvs) == 0 {
		return
	}
	var (
		noExistKVs = make([]KV, 0)
	)

	for _, kv := range kvs {
		_, isExists := m.hashmap.Load(kv.Key)
		if !isExists {
			noExistKVs = append(noExistKVs, kv)
		} else {
			logger.Warn("hashmap exist key", "key", kv.Key)
		}
	}
	if len(noExistKVs) == 0 {
		return
	}

	m.linklistLock.Lock()
	defer m.linklistLock.Unlock()
	for _, kv := range noExistKVs {
		linkListNode := m.linklist.PushBack(kv.Key)
		m.hashmap.Store(kv.Key, &Node{
			node: linkListNode,
			val:  kv.Val,
		})
	}
	return
}

// GetLinkList 获取linklist
func (m *LinkedHashMap) GetLinkList() *list.List {
	return m.linklist
}
