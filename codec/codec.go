// description: chain5j 
// 
// @author: xwc1125
// @date: 2020/10/19
package codec

import (
	"github.com/chain5j/chain5j-pkg/codec/rlp"
)

var coder Codec

func RegisterCodec(codec Codec) {
	coder = codec
}

func Coder() Codec {
	if coder == nil {
		coder = rlp.NewCodec()
	}
	return coder
}

type Codec interface {
	Encoder
	Decoder
}

// 编码器
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

// 解密器
type Decoder interface {
	Decode(data []byte, structPrt interface{}) error
}
