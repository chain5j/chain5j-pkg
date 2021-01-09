// description: chain5j 
// 
// @author: xwc1125
// @date: 2021/1/5
package codec

import "sync"

var (
	lock  sync.RWMutex
	coder Codec
)

func RegisterCodec(_codec Codec) {
	lock.Lock()
	defer lock.Unlock()
	if coder == nil {
		coder = _codec
	}
}

func Codecor() Codec {
	if coder == nil {
		RegisterCodec(DefaultCodec)
	}
	return coder
}
