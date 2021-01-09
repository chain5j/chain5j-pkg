// description: chain5j-pkg
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

func (this *LinkedHashMap) Lock() {
	this.mutex.Lock()
}

func (this *LinkedHashMap) Unlock() {
	this.mutex.Unlock()
}

func (this *LinkedHashMap) RLock() {
	this.mutex.RLock()
}

func (this *LinkedHashMap) RUnlock() {
	this.mutex.RUnlock()
}

// 添加
func (this *LinkedHashMap) Add(key string, val interface{}) bool {
	this.RLock()
	_, isExists := this.hashmap[key]
	this.RUnlock()
	if isExists {
		return false
	}

	this.Lock()
	defer this.Unlock()
	linkListNode := this.linklist.AddToTail(key)
	this.hashmap[key] = &LinkedHashMapNode{
		linklistNode: linkListNode,
		val:          val,
	}

	return true
}

// 获取数据
func (this *LinkedHashMap) Get(key string) interface{} {
	this.RLock()
	originLinkedHashMapNode, isExists := this.hashmap[key]
	this.RUnlock()
	if !isExists {
		return nil
	}

	return (originLinkedHashMapNode.(*LinkedHashMapNode)).val
}

// 删除
func (this *LinkedHashMap) Remove(key string) (bool, interface{}) {
	this.Lock()
	originLinkedHashMapNode, isExists := this.hashmap[key]
	this.Unlock()
	if !isExists {
		i := this.Len()
		_ = i
		return false, nil
	}

	linkedHashMapNode := originLinkedHashMapNode.(*LinkedHashMapNode)

	this.Lock()
	delete(this.hashmap, key)
	this.Unlock()
	this.linklist.RemoveNode(linkedHashMapNode.linklistNode)
	return true, linkedHashMapNode.val
}

// 获取len
func (this *LinkedHashMap) Len() int {
	return len(this.hashmap)
}

// =========批量处理=========
type KV struct {
	Key string
	Val interface{}
}

// 批量添加
func (this *LinkedHashMap) BatchAdd(kvs ...KV) []KV {
	if kvs == nil || len(kvs) == 0 {
		return nil
	}
	var (
		existKVs   = make([]KV, 0)
		noExistKVs = make([]KV, 0)
	)

	this.RLock()
	for _, kv := range kvs {
		_, isExists := this.hashmap[kv.Key]
		if isExists {
			existKVs = append(existKVs, kv)
		} else {
			noExistKVs = append(noExistKVs, kv)
		}
	}
	this.RUnlock()
	if len(noExistKVs) == 0 {
		return existKVs
	}

	this.Lock()
	defer this.Unlock()
	for _, kv := range noExistKVs {
		linkListNode := this.linklist.AddToTail(kv.Key)
		this.hashmap[kv.Key] = &LinkedHashMapNode{
			linklistNode: linkListNode,
			val:          kv.Val,
		}
	}

	return existKVs
}

// 批量获取
func (this *LinkedHashMap) BatchGet(keys ...string) []KV {
	if keys == nil || len(keys) == 0 {
		return nil
	}
	this.RLock()
	kvs := make([]KV, len(keys))
	for i, key := range keys {
		originLinkedHashMapNode, isExists := this.hashmap[key]
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
	this.RUnlock()
	return kvs
}

// 批量删除
func (this *LinkedHashMap) BatchRemove(keys ...string) []KV {
	if keys == nil || len(keys) == 0 {
		return nil
	}
	kvs := make([]KV, 0)
	this.RLock()
	for _, key := range keys {
		originLinkedHashMapNode, isExists := this.hashmap[key]
		if isExists {
			linkedHashMapNode := originLinkedHashMapNode.(*LinkedHashMapNode)
			kvs = append(kvs, KV{
				Key: key,
				Val: linkedHashMapNode,
			})
		}
	}
	this.RUnlock()

	this.Lock()
	for _, kv := range kvs {
		delete(this.hashmap, kv.Key)
		this.linklist.RemoveNode(kv.Val.(*LinkedHashMapNode).linklistNode)
	}
	this.Unlock()
	return kvs
}

func (this *LinkedHashMap) GetLinkList() *LinkList {
	return this.linklist
}
