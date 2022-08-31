// Package treemap
//
// @author: xwc1125
package treemap

import (
	"encoding/json"
	"io"

	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/collection/maps/treemap/v2"
	"golang.org/x/exp/constraints"
)

// TreeMap based on red-black tree, alias of RedBlackTree.
type TreeMap[K constraints.Ordered, V any] struct {
	*treemap.TreeMap[K, V]
}

func NewKvs[K constraints.Ordered, V any](t *treemap.TreeMap[K, V]) *TreeMap[K, V] {
	return &TreeMap[K, V]{
		TreeMap: t,
	}
}

func (k TreeMap[K, V]) EncodeRLP(w io.Writer) error {
	// [xwc1125] 如果使用，encoding/json返回的是64进制
	// 如果使用，github.com/json-iterator/go返回的是16进制
	json, err := MarshalJSON(k.TreeMap)
	if err != nil {
		return err
	}
	return rlp.Encode(w, json)
}

func (k *TreeMap[K, V]) DecodeRLP(s *rlp.Stream) error {
	var jsonBytes []byte
	err := s.Decode(&jsonBytes)
	if err != nil {
		return err
	}
	var treeMap1 = treemap.New[K, V]()
	err = UnmarshalJSON(jsonBytes, treeMap1)
	if err != nil {
		return err
	}
	*k.TreeMap = *treeMap1
	return err
}

func (k *TreeMap[K, V]) Clone() *TreeMap[K, V] {
	clone := Clone(k.TreeMap)
	return &TreeMap[K, V]{
		TreeMap: clone,
	}
}

func (k *TreeMap[K, V]) MarshalJSON() ([]byte, error) {
	return MarshalJSON(k.TreeMap)
}

func MarshalJSON[K constraints.Ordered, V any](t *treemap.TreeMap[K, V]) ([]byte, error) {
	return json.Marshal(Map(t))
}

func Map[K constraints.Ordered, V any](t *treemap.TreeMap[K, V]) map[K]V {
	m := make(map[K]V, t.Len())
	for it := t.Iterator(); it.Valid(); it.Next() {
		m[it.Key()] = it.Value()
	}
	return m
}

func (k *TreeMap[K, V]) UnmarshalJSON(data []byte) error {
	if k.TreeMap == nil {
		k.TreeMap = treemap.New[K, V]()
	}
	return UnmarshalJSON(data, k.TreeMap)
}

func UnmarshalJSON[K constraints.Ordered, V any](b []byte, t *treemap.TreeMap[K, V]) error {
	var data map[K]V
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for k, v := range data {
		t.Set(k, v)
	}
	return nil
}

func Clone[K constraints.Ordered, V any](t *treemap.TreeMap[K, V]) *treemap.TreeMap[K, V] {
	newTree := treemap.NewWithKeyCompare[K, V](t.KeyCompare())
	m := Map(t)
	for k, v := range m {
		newTree.Set(k, v)
	}
	return newTree
}
