// Package linkedHashMap
//
// @author: xwc1125
// @date: 2020/12/22
package linkedHashMap

import (
	"sync"
)

type LinkedHashMapNode struct {
	linklistNode *LinkListNode
	val          interface{}
}

type LinkedHashMap struct {
	linklist *LinkList
	hashmap  map[string]interface{}
	mutex    *sync.RWMutex
}

func NewLinkedHashMap() *LinkedHashMap {
	return &LinkedHashMap{
		linklist: NewLinkList(),
		hashmap:  make(map[string]interface{}),
		mutex:    &sync.RWMutex{},
	}
}

func (hm *LinkedHashMap) Lock() {
	hm.mutex.Lock()
}

func (hm *LinkedHashMap) Unlock() {
	hm.mutex.Unlock()
}

func (hm *LinkedHashMap) RLock() {
	hm.mutex.RLock()
}

func (hm *LinkedHashMap) RUnlock() {
	hm.mutex.RUnlock()
}

// Add 添加
func (hm *LinkedHashMap) Add(key string, val interface{}) bool {
	hm.RLock()
	_, isExists := hm.hashmap[key]
	hm.RUnlock()
	if isExists {
		return false
	}

	hm.Lock()
	defer hm.Unlock()
	linkListNode := hm.linklist.AddToTail(key)
	hm.hashmap[key] = &LinkedHashMapNode{
		linklistNode: linkListNode,
		val:          val,
	}

	return true
}

// Get 获取数据
func (hm *LinkedHashMap) Get(key string) interface{} {
	hm.RLock()
	originLinkedHashMapNode, isExists := hm.hashmap[key]
	hm.RUnlock()
	if !isExists {
		return nil
	}

	return (originLinkedHashMapNode.(*LinkedHashMapNode)).val
}

// Remove 删除
func (hm *LinkedHashMap) Remove(key string) (bool, interface{}) {
	hm.Lock()
	originLinkedHashMapNode, isExists := hm.hashmap[key]
	hm.Unlock()
	if !isExists {
		i := hm.Len()
		_ = i
		return false, nil
	}

	linkedHashMapNode := originLinkedHashMapNode.(*LinkedHashMapNode)

	hm.Lock()
	delete(hm.hashmap, key)
	hm.Unlock()
	hm.linklist.RemoveNode(linkedHashMapNode.linklistNode)
	return true, linkedHashMapNode.val
}

// Len 获取len
func (hm *LinkedHashMap) Len() int {
	return len(hm.hashmap)
}

// =========批量处理=========

// BatchAdd 批量添加
func (hm *LinkedHashMap) BatchAdd(kvs ...KV) []KV {
	if kvs == nil || len(kvs) == 0 {
		return nil
	}
	var (
		existKVs   = make([]KV, 0)
		noExistKVs = make([]KV, 0)
	)

	hm.RLock()
	for _, kv := range kvs {
		_, isExists := hm.hashmap[kv.Key]
		if isExists {
			existKVs = append(existKVs, kv)
		} else {
			noExistKVs = append(noExistKVs, kv)
		}
	}
	hm.RUnlock()
	if len(noExistKVs) == 0 {
		return existKVs
	}

	hm.Lock()
	defer hm.Unlock()
	for _, kv := range noExistKVs {
		linkListNode := hm.linklist.AddToTail(kv.Key)
		hm.hashmap[kv.Key] = &LinkedHashMapNode{
			linklistNode: linkListNode,
			val:          kv.Val,
		}
	}

	return existKVs
}

// BatchGet 批量获取
func (hm *LinkedHashMap) BatchGet(keys ...string) []KV {
	if keys == nil || len(keys) == 0 {
		return nil
	}
	hm.RLock()
	kvs := make([]KV, len(keys))
	for i, key := range keys {
		originLinkedHashMapNode, isExists := hm.hashmap[key]
		if !isExists {
			kvs[i] = KV{
				Key: key,
				Val: nil,
			}
		} else {
			kvs[i] = KV{
				Key: key,
				Val: (originLinkedHashMapNode.(*LinkedHashMapNode)).val,
			}
		}
	}
	hm.RUnlock()
	return kvs
}

// BatchRemove 批量删除
func (hm *LinkedHashMap) BatchRemove(keys ...string) []KV {
	if keys == nil || len(keys) == 0 {
		return nil
	}
	kvs := make([]KV, 0)
	hm.RLock()
	for _, key := range keys {
		originLinkedHashMapNode, isExists := hm.hashmap[key]
		if isExists {
			linkedHashMapNode := originLinkedHashMapNode.(*LinkedHashMapNode)
			kvs = append(kvs, KV{
				Key: key,
				Val: linkedHashMapNode,
			})
		}
	}
	hm.RUnlock()

	hm.Lock()
	for _, kv := range kvs {
		delete(hm.hashmap, kv.Key)
		hm.linklist.RemoveNode(kv.Val.(*LinkedHashMapNode).linklistNode)
	}
	hm.Unlock()
	return kvs
}

// GetLinkList 获取linklist
func (hm *LinkedHashMap) GetLinkList() *LinkList {
	return hm.linklist
}
