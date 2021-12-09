// Package hashmap
//
// @author: xwc1125
package hashmap

import (
	"encoding/json"
	"fmt"
	"github.com/chain5j/chain5j-pkg/codec"
	"github.com/chain5j/chain5j-pkg/util/reflectutil"
	"reflect"
	"sort"
	"sync"
)

// HashMap ...
type HashMap struct {
	data   map[string]interface{} // 数据
	val    interface{}            // 同类型的map类型
	m      sync.Mutex             // 锁
	isSafe bool                   // 是否安全
}

func (m *HashMap) lock() {
	if m.isSafe {
		m.m.Lock()
	}
}

func (m *HashMap) unlock() {
	if m.isSafe {
		m.m.Unlock()
	}
}

// NewHashMap ...
func NewHashMap(isSafe bool) *HashMap {
	return &HashMap{
		data:   make(map[string]interface{}),
		isSafe: isSafe,
	}
}

func NewHashMapSame(isSafe bool, val interface{}) *HashMap {
	h := NewHashMap(isSafe)
	if val != nil {
		h.val = val
	}
	return h
}

func NewHashMapFill(val interface{}) *HashMap {
	valueOf := reflect.ValueOf(val)
	if valueOf.Kind() == reflect.Map {
		//elem := valueOf.Elem()
		//fmt.Println(elem)
		_, err := Get(val, []interface{}{})
		if err != nil {
			fmt.Println(err)
		}
		m := valueOf.Interface().(map[string]interface{})
		fmt.Println(m)
	}
	fmt.Println(valueOf.Kind().String())
	h := &HashMap{
		data: val.(map[string]interface{}),
	}
	return h
}

func GetMapI(v interface{}, path ...interface{}) (map[interface{}]interface{}, error) {
	v, err := Get(v, path...)
	if err != nil {
		return nil, err
	}
	m, ok := v.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map with interface keys node, got: %T", v)
	}
	return m, nil
}

// GetMapS returns a map with string keys denoted by the path.
//
// If path is empty or nil, v is returned as a slice.
func GetMapS(v interface{}, path ...interface{}) (map[string]interface{}, error) {
	v, err := Get(v, path...)
	if err != nil {
		return nil, err
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map with string keys node, got: %T", v)
	}
	return m, nil
}

func Get(v interface{}, path ...interface{}) (interface{}, error) {
	for i, el := range path {
		switch node := v.(type) {
		case map[string]interface{}:
			key, ok := el.(string)
			if !ok {
				return nil, fmt.Errorf("expected string path element, got: %T (path element idx: %d)", el, i)
			}
			v, ok = node[key]
			if !ok {
				return nil, fmt.Errorf("missing key: %s (path element idx: %d)", key, i)
			}

		case map[interface{}]interface{}:
			var ok bool
			v, ok = node[el]
			if !ok {
				return nil, fmt.Errorf("missing key: %v (path element idx: %d)", el, i)
			}

		case []interface{}:
			idx, ok := el.(int)
			if !ok {
				return nil, fmt.Errorf("expected int path element, got: %T (path element idx: %d)", el, i)
			}
			if idx < 0 || idx >= len(node) {
				return nil, fmt.Errorf("index out of range: %d (path element idx: %d)", idx, i)
			}
			v = node[idx]

		default:
			return nil, fmt.Errorf("expected map or slice node, got: %T (path element idx: %d)", node, i)
		}
	}

	return v, nil
}

// SetSafe 设置为是否安全
func (m *HashMap) SetSafe(isSafe bool) {
	m.lock()
	m.isSafe = isSafe
	m.unlock()
}

// Put 增加或者修改一个元素
func (m *HashMap) Put(k string, v interface{}) {
	m.lock()
	m.data[k] = v
	m.unlock()
}

// Get ...
// @obj: 值对象
// @objType: 值类型
// @isExist: 是否存在
func (m *HashMap) Get(k string) (obj interface{}, objType string, isExist bool) {
	m.lock()
	v, ok := m.data[k]
	m.unlock()
	var rv interface{} = nil
	var rt = ""
	if ok {
		rv = v
		rt = reflect.TypeOf(v).String()
	}
	return rv, rt, ok
}

// GetObj ...
func (m *HashMap) GetObj(k string) (obj interface{}) {
	m.lock()
	v, ok := m.data[k]
	m.unlock()
	var rv interface{} = nil
	if ok {
		rv = v
	}
	return rv
}

// GetValue ...
func (m *HashMap) GetValue(k string) reflect.Value {
	m.lock()
	v, ok := m.data[k]
	m.unlock()
	var rv interface{} = nil
	var sv reflect.Value
	if ok {
		rv = v
		sv = reflectutil.GetValue(rv)
	}
	return sv
}

// GetValues ...
func (m *HashMap) GetValues(k string) []reflect.Value {
	m.lock()
	v, ok := m.data[k]
	m.unlock()
	var rv interface{} = nil
	var sv []reflect.Value
	if ok {
		rv = v
		sv = reflectutil.GetValues(rv)
	}
	return sv
}

// HasKey 判断是否包括key，如果包含key返回value的类型
func (m *HashMap) HasKey(k string) (bool, string) {
	m.lock()
	v, ok := m.data[k]
	m.unlock()
	var rt = ""
	if ok {
		rt = reflect.TypeOf(v).String()
	}
	return ok, rt
}

// Exist ...
func (m *HashMap) Exist(k string) bool {
	m.lock()
	_, ok := m.data[k]
	m.unlock()
	return ok
}

// Remove 移除一个元素
func (m *HashMap) Remove(k string) (interface{}, bool) {
	m.lock()
	v, ok := m.data[k]
	var rv interface{} = nil
	if ok {
		rv = v
		delete(m.data, k)
	}
	m.unlock()
	return rv, ok
}

// ForEach 复制map用于外部遍历
func (m *HashMap) ForEach() map[string]interface{} {
	m.lock()
	mb := map[string]interface{}{}
	for k, v := range m.data {
		mb[k] = v
	}
	m.unlock()
	return mb
}

// Size ...
func (m *HashMap) Size() int {
	return len(m.data)
}

// Sort 排序
func (m *HashMap) Sort() []KV {
	var kvs []KV
	var keyArray []string
	m.lock()
	for k, _ := range m.data {
		keyArray = append(keyArray, k)
	}
	sort.Strings(keyArray)
	for _, v := range keyArray {
		kv := &KV{
			Key:   v,
			Value: m.data[v],
		}
		kvs = append(kvs, *kv)
	}
	m.unlock()
	return kvs
}

// String ...
func (m *HashMap) String() string {
	bytes, _ := json.Marshal(m.data)
	return string(bytes)
}

func (m *HashMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

func (m *HashMap) UnmarshalJSON(bytes []byte) error {
	if m.data == nil {
		m.data = make(map[string]interface{})
		m.isSafe = true
	}
	return json.Unmarshal(bytes, &m.data)
}

func (m *HashMap) Decode(data []byte) error {
	return codec.Coder().Decode(data, &m)
}
