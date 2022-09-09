// Package hashmap
//
// @author: xwc1125
package hashmap

import (
	"io"
	"math/big"
	"reflect"

	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/convutil"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/mitchellh/mapstructure"
)

func (m *HashMap) EncodeRLP(w io.Writer) error {
	data := m.Sort()
	for i, kv := range data {
		switch vv := kv.Value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			data[i] = KV{
				Key:   kv.Key,
				Value: convutil.ToUint64(vv),
			}
		}
	}
	return rlp.Encode(w, data)
}

func (m *HashMap) DecodeRLP(s *rlp.Stream) error {
	var data []KV
	if err := s.Decode(&data); err != nil {
		return err
	}
	for _, kv := range data {
		m.data[kv.Key] = kv.Value
	}
	return nil
}

// ToStruct 用map填充结构
func (m *HashMap) ToStruct(out interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       StringToByteSizesHookFunc,
		Metadata:         nil,
		Result:           out,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(m.data)
}

func StringToByteSizesHookFunc(
	f reflect.Type,
	des reflect.Type,
	data interface{}) (interface{}, error) {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	// dataType := dataVal.Type()
	desKind := getKindByKind(des.Kind())

	switch {
	case dataKind == reflect.String:
		if i, ok := parseStr(des, data.(string)); ok {
			return i, nil
		}
	case dataKind == reflect.Slice,
		dataKind == reflect.Array:
		// todo xwc1125 []byte转换
		dataType := dataVal.Type()
		elemKind := dataType.Elem().Kind()
		switch elemKind {
		case reflect.Uint8:
			var uints []uint8
			if dataKind == reflect.Array {
				uints = make([]uint8, dataVal.Len(), dataVal.Len())
				for i := range uints {
					uints[i] = dataVal.Index(i).Interface().(uint8)
				}
			} else {
				uints = dataVal.Interface().([]uint8)
			}
			switch {
			case desKind == reflect.Uint:
				return convutil.BytesToUint64(uints), nil
			case desKind == reflect.Int:
				return convutil.BytesToInt64(uints), nil
			case desKind == reflect.Bool:
				return convutil.BytesToBool(uints)
			default:
				if i, ok := parseBytes(des, uints); ok {
					return i, nil
				}
			}
		}
	}
	return data, nil
}

func parseStr(des reflect.Type, data string) (interface{}, bool) {
	switch des {
	case reflect.TypeOf(types.Address{}):
		address := types.HexToAddress(data)
		return address, true
	case reflect.TypeOf(types.Hash{}):
		hash := types.HexToHash(data)
		return hash, true
	case reflect.TypeOf(math.HexOrDecimal256{}):
		result := new(math.HexOrDecimal256)
		result.UnmarshalText([]byte(data))
		return result, true
	case reflect.TypeOf(big.Int{}):
		result := new(big.Int)
		result.UnmarshalText([]byte(data))
		return result, true
	case reflect.TypeOf(hexutil.Bytes{}):
		result := hexutil.MustDecode(data)
		return result, true
	}
	return data, false
}
func parseBytes(des reflect.Type, data []byte) (interface{}, bool) {
	switch des {
	case reflect.TypeOf(types.Address{}):
		if len(data) == types.AddressLength {
			return types.BytesToAddress(data), true
		}
		return types.HexToAddress(string(data)), true
	case reflect.TypeOf(types.Hash{}):
		if len(data) == types.HashLength {
			return types.BytesToHash(data), true
		}
		return types.HexToHash(string(data)), true
	case reflect.TypeOf(math.HexOrDecimal256{}):
		result := new(math.HexOrDecimal256)
		result.UnmarshalText(data)
		return result, true
	case reflect.TypeOf(big.Int{}):
		result := new(big.Int)
		result.UnmarshalText(data)
		return result, true
	case reflect.TypeOf(hexutil.Bytes{}):
		return data, true
	}
	return data, false
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()
	return getKindByKind(kind)
}
func getKindByKind(kind reflect.Kind) reflect.Kind {
	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}
