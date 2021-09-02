// Package rlp
// 
// @author: xwc1125
// @date: 2020/10/20
package rlp

type Codec struct {
}

func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) Encode(v interface{}) ([]byte, error) {
	return EncodeToBytes(v)
}

func (c *Codec) Decode(data []byte, structPrt interface{}) error {
	return DecodeBytes(data, structPrt)
}
