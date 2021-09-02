// Package codec
//
// @author: xwc1125
// @date: 2021/1/5
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

// Coder 获取全局的编解码器
func Coder() Codec {
	if coder == nil {
		return DefaultCodec
	}
	return coder
}
