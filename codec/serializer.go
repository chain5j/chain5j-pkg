// Package codec
//
// @author: xwc1125
// @date: 2021/8/28
package codec

// Serializer 序列化
type Serializer interface {
	Serialize() ([]byte, error)
}

// Deserializer 反序列化
type Deserializer interface {
	Deserialize(d []byte) error
}
