// description: chain5j-devops-tools 
// 
// @author: xwc1125
// @date: 2020/5/29 0029
package rpc

import "context"

type JsonRpc interface {
	Call(result interface{}, method string, args ...interface{}) error
	CallContext(ctx context.Context, result interface{}, method string, args1 ...interface{}) error
	Close()
}
