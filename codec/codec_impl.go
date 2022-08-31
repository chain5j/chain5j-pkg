// Package codec
//
// @author: xwc1125
package codec

import "sync"

var (
	once  sync.Once
	coder Codec
)

// RegisterCodec 注册编解码器
func RegisterCodec(_codec Codec) {
	once.Do(func() {
		if coder == nil {
			coder = _codec
		}
	})
}

// Codecor 获取全局的编解码器
func Codecor() Codec {
	if coder == nil {
		return DefaultCodec
	}
	return coder
}
