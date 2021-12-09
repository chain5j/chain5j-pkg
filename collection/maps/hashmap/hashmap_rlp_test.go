// Package hashmap
//
// @author: xwc1125
package hashmap

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/codec"
	"github.com/chain5j/chain5j-pkg/math"
	"github.com/chain5j/chain5j-pkg/util/convutil"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/mitchellh/mapstructure"
	"testing"
)

func TestConvert(t *testing.T) {
	//int8: -128 ~ 127
	//int16: -32768 ~ 32767
	//int32: -2147483648 ~ 2147483647
	//int64: -9223372036854775808 ~ 9223372036854775807
	//uint8: 0 ~ 255
	//uint16: 0 ~ 65535
	//uint32: 0 ~ 4294967295
	//uint64: 0 ~ 18446744073709551615
	hashMap := NewHashMap(true)
	hashMap.Put("int_0", 0)
	hashMap.Put("int_1", 1)
	hashMap.Put("int8_MinInt8", math.MinInt8)
	hashMap.Put("int8_MaxInt8", math.MaxInt8)
	hashMap.Put("int16_MinInt16", math.MinInt16)
	hashMap.Put("int16_MaxInt16", math.MaxInt16)
	hashMap.Put("int32_MinInt32", math.MinInt32)
	hashMap.Put("int32_MaxInt32", math.MaxInt32)
	hashMap.Put("int64_MinInt64", math.MinInt64)
	hashMap.Put("int64_MaxInt64", math.MaxInt64)
	// uint
	hashMap.Put("uint_0", 0)
	hashMap.Put("uint8_MaxUint8", math.MaxUint8)
	hashMap.Put("uint16_MaxUint16", math.MaxUint16)
	hashMap.Put("uint32_MaxUint32", math.MaxUint32)
	//hashMap.Put("uint64_1", math.MaxUint64)

	hashMap.Put("show", true)

	if bytes, err := codec.Coder().Encode(hashMap); err != nil {
		t.Fatal(err)
	} else {
		newHashMap := NewHashMap(true)
		if err := codec.Coder().Decode(bytes, newHashMap); err != nil {
			t.Fatal(err)
		}
		for k, v := range newHashMap.data {
			t.Log("---------------" + k + "-------------")
			{
				t.Log("===========BytesToBool============")
				b, err := convutil.BytesToBool(v.([]byte))
				if err != nil {
					t.Error(k, err)
				} else {
					t.Log(k, b)
				}
			}
			{
				t.Log("===========BytesToInt64============")
				b, err := convutil.BytesToInt64(v.([]byte))
				if err != nil {
					t.Error(k, err)
				} else {
					t.Log(k, b)
				}
			}
			{
				t.Log("===========BytesToInt64U============")
				b, err := convutil.BytesToUint64(v.([]byte))
				if err != nil {
					t.Error(k, err)
				} else {
					t.Log(k, b)
				}
			}
			{
				t.Log("===========BytesToString============")
				b, err := convutil.BytesToString(v.([]byte))
				if err != nil {
					t.Error(k, err)
				} else {
					t.Log(k, b)
				}
			}
		}
	}
}

func TestRlp(t *testing.T) {
	hashMap := NewHashMap(true)
	hashMap.Put("epoch", 30000)
	hashMap.Put("period", 5)
	hashMap.Put("policy", 0)
	hashMap.Put("timeout", 10000)
	hashMap.Put("show", true)

	hashMap.Put("manager", "0x9254e62fbca63769dfd4cc8e23f630f0785610ce0x9254e62fbca63769dfd4cc8e23f630f0785610ce")
	hashMap.Put("managers", []string{
		"0x9254e62fbca63769dfd4cc8e23f630f0785610ce0x9254e62fbca63769dfd4cc8e23f630f0785610ce",
		"0x92c8cae42a94045670cbb0bfcf8f790d9f8097e70x9254e62fbca63769dfd4cc8e23f630f0785610ce",
	})
	hashMap.Put("validator", []string{
		"QmVy5JASWLUns3Wwe91arP93rPPumkEQ7fQZ4GfpevxKbd",
		"QmafFDAXGtW1M2zhsiYPxUtpkYjAgBUpPDPQZRkhgonGrT",
	})
	if bytes, err := codec.Coder().Encode(hashMap); err != nil {
		t.Fatal(err)
	} else {
		newHashMap := NewHashMap(true)
		if err := codec.Coder().Decode(bytes, newHashMap); err != nil {
			t.Fatal(err)
		}
		t.Logf("%v", newHashMap)
		{
			obj := newHashMap.GetObj("managers")
			switch p := obj.(type) {
			case []interface{}:
				for _, p1 := range p {
					v := p1.([]byte)
					t.Log(string(v))
				}
			}
		}
		{
			obj := newHashMap.GetObj("epoch")
			switch p := obj.(type) {
			case []byte:
				toInt, _ := convutil.BytesToInt64(p)
				fmt.Println(toInt)
				fmt.Println(hexutil.Encode(p))
			}
		}
		for k, v := range newHashMap.data {
			switch k {
			case "show":
				b, err := convutil.BytesToBool(v.([]byte))
				if err != nil {
					t.Error(k, err)
				} else {
					t.Log(k, b)
				}
			case "manager":
				b, err := convutil.BytesToString(v.([]byte))
				if err != nil {
					t.Error(k, err)
				} else {
					t.Log(k, b)
				}
			case "epoch", "epoch1":
				b, err := convutil.BytesToUint64(v.([]byte))
				if err != nil {
					t.Error(k, err)
				} else {
					t.Log(k, b)
				}
			default:
				switch v.(type) {
				case []byte:
					b, err := convutil.BytesToInt64(v.([]byte))
					if err != nil {
						t.Error(k, err)
					} else {
						t.Log(k, b)
					}
				}
			}
		}

		var config testConfig
		if err := mapstructure.WeakDecode(newHashMap.data, &config); err != nil {
			t.Error(err)
		} else {
			t.Log(config)
		}
	}
}

// testConfig 配置
type testConfig struct {
	Epoch     uint64   `json:"epoch" mapstructure:"epoch"`         // 期数
	Timeout   uint64   `json:"timeout" mapstructure:"timeout"`     // round 超时时间， 在BlockPeriod不为0时有效
	Period    uint64   `json:"period" mapstructure:"period"`       // 区块产生间隔
	Managers  []string `json:"managers" mapstructure:"managers"`   // 管理员
	Validator []string `json:"validator" mapstructure:"validator"` // 验证者
	Show      bool     `json:"show" mapstructure:"show"`
}
