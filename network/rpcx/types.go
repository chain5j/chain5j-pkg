// Package rpcx
//
// @author: xwc1125
package rpcx

import (
	"encoding/json"
	"reflect"
	"sync"
)

// service 服务
type service struct {
	name      string       // 服务名称
	typ       reflect.Type // 接收者类型
	callbacks sync.Map     // 注册回调函数,map[string]*callback
}

type callback struct {
	rcvr     reflect.Value  // 方法接收者
	method   reflect.Method // callback
	argTypes []reflect.Type // 参数类型
	hasCtx   bool           // 方法的第一个参数是否为context
	errPos   int            // 错误的位置，如果无错误=-1
}
type callbacks map[string]*callback

type serverRequest struct {
	id      interface{}
	svcname string
	callb   *callback
	args    []reflect.Value
	err     error
}
type Error interface {
	Error() string  // returns the message
	ErrorCode() int // returns the code
}

type JsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (j JsonError) String() string {
	marshal, _ := json.Marshal(j)
	return string(marshal)
}

type JsonResponse struct {
	Version string      `json:"jsonrpc,omitempty"`
	Id      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   error       `json:"error,omitempty"`
}
type jsonResponseExt struct {
	Version string      `json:"jsonrpc,omitempty"`
	Id      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (j JsonResponse) String() string {
	marshal, _ := json.Marshal(j)
	return string(marshal)
}

func (j JsonResponse) MarshalJSON() ([]byte, error) {
	ext := jsonResponseExt{
		Version: j.Version,
		Id:      j.Id,
		Result:  j.Result,
	}
	if j.Error != nil {
		ext.Error = j.Error.Error()
	}
	return json.Marshal(ext)
}

type Namespace struct {
	Namespace string    `json:"namespace"`
	Methods   []*Method `json:"methods"`
}
type Method struct {
	Method   string   `json:"method" mapstructure:"method" yaml:"method"`
	ArgsType []string `json:"argsType,omitempty" mapstructure:"method,omitempty" yaml:"method,omitempty"`
}

type API struct {
	Namespace     string      // 对外的命名空间
	Version       string      // api版本
	Public        bool        // 可供公众使用
	Service       interface{} // 方法集合
	Authenticated bool        // 是否提供认证后才能访问
}

type Header map[string][]string
