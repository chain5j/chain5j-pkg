// Package codec
//
// @author: xwc1125
// @date: 2020/10/19
package codec

import "github.com/chain5j/chain5j-pkg/codec/rlp"

var DefaultCodec = rlp.NewCodec()

// Codec 编解码
type Codec interface {
	Encoder
	Decoder
}

// Encoder 编码器
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

// Decoder 解码器
type Decoder interface {
	Decode(data []byte, structPrt interface{}) error
}
