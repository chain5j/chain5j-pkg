// description: chain5j 
// 
// @author: xwc1125
// @date: 2020/11/15
package hashmap

import (
	"encoding/json"
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/util/reflectutil"
	"log"
	"reflect"
)

type KVBytes struct {
	Key   string
	Value []byte
}

type Encoder interface {
	EncodeToBytes() ([]byte, error)
}

type Decoder interface {
	DecodeFromBytes([]byte) error
}

var (
	EncoderInterface = reflect.TypeOf(new(Encoder)).Elem()
	DecoderInterface = reflect.TypeOf(new(Decoder)).Elem()
)

func (m *HashMap) EncodeToBytes() ([]byte, error) {
	var data []KVBytes
	var kv KVBytes
	sort2 := m.Sort() // 排序

	for _, v := range sort2 {
		kv = KVBytes{
			Key: v.Key,
		}
		v1 := reflectutil.ToPointer(v.Value)
		valueOf := reflect.ValueOf(v1)
		if valueOf.Type().Implements(EncoderInterface) {
			bytes, err := valueOf.Interface().(Encoder).EncodeToBytes()
			if err != nil {
				log.Println(err)
				return nil, err
			}
			kv.Value = bytes
		} else {
			bytes, err := json.Marshal(v.Value)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			kv.Value = bytes
		}
		data = append(data, kv)
	}
	return rlp.EncodeToBytes(data)
}

func (m *HashMap) DecodeFromBytes(b []byte) error {
	var data []KVBytes
	if err := rlp.DecodeBytes(b, &data); err != nil {
		return err
	}
	var valueOf reflect.Value
	val := m.val
	if m.val != nil {
		val = reflectutil.ToPointer(m.val)
		valueOf = reflect.ValueOf(val)
	}
	for _, v := range data {
		if val != nil {
			if valueOf.Type().Implements(DecoderInterface) {
				val := valueOf.Interface().(Decoder)
				err := val.DecodeFromBytes(v.Value)
				if err != nil {
					return err
				}
			} else {
				err := json.Unmarshal(v.Value, val)
				if err != nil {
					return err
				}
			}
			if reflect.TypeOf(m.val).Kind() != reflect.Ptr {
				m.data[v.Key] = reflectutil.DelPointer(val)
			} else {
				m.data[v.Key] = val
			}
		} else {
			m.data[v.Key] = v.Value
		}
	}
	return nil
}

func (m *HashMap) DecodeFromBytesMult(b []byte, vals []interface{}) error {
	var data []KVBytes
	if err := rlp.DecodeBytes(b, &data); err != nil {
		return err
	}
	for _, v := range data {
		if vals != nil {
			for _, val := range vals {
				if val != nil {
					val = reflectutil.ToPointer(val)
					valueOf := reflect.ValueOf(val)
					if valueOf.Type().Implements(DecoderInterface) {
						val := valueOf.Interface().(Decoder)
						err := val.DecodeFromBytes(v.Value)
						if err != nil {
							continue
						}
					} else {
						err := json.Unmarshal(v.Value, val)
						if err != nil {
							continue
						}
					}
					if reflect.TypeOf(val).Kind() != reflect.Ptr {
						m.data[v.Key] = reflectutil.DelPointer(val)
					} else {
						m.data[v.Key] = val
					}
				}
			}
		} else {
			m.data[v.Key] = v.Value
		}
	}
	return nil
}
