// Package json
//
// @author: xwc1125
package json

import (
	"github.com/chain5j/chain5j-pkg/codec"
)

var _ codec.Codec = &Codec{}

type Codec struct {
}

func NewCodec() *Codec {
	return &Codec{}
}

func (c *Codec) Encode(v interface{}) ([]byte, error) {
	return Marshal(v)
}

func (c *Codec) Decode(data []byte, structPrt interface{}) error {
	return Unmarshal(data, structPrt)
}
