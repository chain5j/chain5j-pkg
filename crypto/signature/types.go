// Package signature
//
// @author: xwc1125
package signature

import (
	"github.com/chain5j/chain5j-pkg/codec"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

// SignResult 签名结果
type SignResult struct {
	Name      string        `json:"name" mapstructure:"name"`           // 算法名称
	PubKey    hexutil.Bytes `json:"pub_key" mapstructure:"pub_key"`     // 公钥
	Signature hexutil.Bytes `json:"signature" mapstructure:"signature"` // 签名结果
}

// Serialize 获取signResult的bytes
func (s *SignResult) Serialize() ([]byte, error) {
	return codec.Codecor().Encode(s)
}

// Deserialize 反解码
func (s *SignResult) Deserialize(data []byte) error {
	var sign SignResult
	err := codec.Codecor().Decode(data, &sign)
	if err != nil {
		return err
	}
	*s = sign
	return nil
}

func (s *SignResult) Copy() *SignResult {
	cpy := SignResult{
		Name: s.Name,
	}
	if !s.PubKey.Nil() {
		cpy.PubKey = make([]byte, len(s.PubKey))
		copy(cpy.PubKey, s.PubKey)
	}
	if len(s.Signature) > 0 {
		cpy.Signature = make([]byte, len(s.Signature))
		copy(cpy.Signature, s.Signature)
	}
	return &cpy
}
