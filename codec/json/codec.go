// description: chain5j 
// 
// @author: xwc1125
// @date: 2020/10/20
package json

import (
	"github.com/chain5j/chain5j-pkg/codec"
	"github.com/chain5j/chain5j-pkg/json"
)

var _ codec.Codec = &Codec{}

type Codec struct {
}

func Newcodec.Codecor() *Codec {
	return &Codec{}
}

func (c *Codec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *Codec) Decode(data []byte, structPrt interface{}) error {
	return json.Unmarshal(data, structPrt)
}
