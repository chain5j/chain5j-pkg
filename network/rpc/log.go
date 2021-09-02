// Package rpc
//
// @author: xwc1125
// @date: 2021/1/9
package rpc

import (
	"github.com/chain5j/chain5j-pkg/logger"
)

func log15() logger.Logger {
	return logger.New("rpc")
}
