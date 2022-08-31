// Package rpcx
//
// @author: xwc1125
package rpcx

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/chain5j/chain5j-pkg/util/convutil"
	"github.com/chain5j/logger"
)

var (
	JsonRpcVersion         = "2.0"
	ServiceMethodSeparator = "."
	Debug                  = false
)

// RPC rpc 的接口
type RPC interface {
	// Register 注册命名及其类方法[对外暴露的方法必须是首字母大写]
	// namespace 命名空间
	// rcvr 服务对象(函数类)
	Register(namespace string, rcvr interface{}) error
	// Exec 执行
	// namespace 命名空间
	// method 方法名称
	// reqId 请求的ID
	// params 请求的参数
	Exec(ctx context.Context, namespace string, method string, reqId interface{}, params ...interface{}) *JsonResponse

	Namespaces(ctx context.Context) []*Namespace
	Namespaces2(ctx context.Context) map[string]map[string]*Method
}

var (
	_ RPC = new(RPCX)
)

// RPCX rpcx
type RPCX struct {
	log      logger.Logger
	services sync.Map // map[string]*service

	namespaces    []*Namespace
	namespacesMap map[string]map[string]*Method
}

// New 创建rpc
func New() RPC {
	return &RPCX{
		log:      logger.Log("rpcx"),
		services: sync.Map{},
	}
}

// Register 注册
func (r *RPCX) Register(namespace string, rcvr interface{}) error {
	namespace = formatName(namespace)
	svc := new(service)
	svc.typ = reflect.TypeOf(rcvr)
	rcvrVal := reflect.ValueOf(rcvr)
	if namespace == "" {
		return fmt.Errorf("no service name for type %s", svc.typ.String())
	}
	if !isExported(reflect.Indirect(rcvrVal).Type().Name()) {
		return fmt.Errorf("%s is not exported", reflect.Indirect(rcvrVal).Type().Name())
	}
	methods := suitableCallbacks(rcvrVal, svc.typ)

	if len(methods) == 0 {
		return fmt.Errorf("service %T doesn't have any suitable methods to expose", rcvr)
	}

	if regsvc, present := r.services.Load(namespace); present {
		s := regsvc.(*service)
		for _, m := range methods {
			s.callbacks.Store(formatName(m.method.Name), m)
		}
		r.services.Store(namespace, s)
		return nil
	}

	svc.name = namespace
	for _, m := range methods {
		svc.callbacks.Store(formatName(m.method.Name), m)
	}

	r.services.Store(svc.name, svc)
	return nil
}

// suitableCallbacks 获取命名空间有哪些方法是暴露的
// rcvr 传入对象的值类型
// typ 传入对象的类型。如果是类，那么就可以获取其所有的方法
func suitableCallbacks(rcvr reflect.Value, typ reflect.Type) callbacks {
	callbacks := make(callbacks)

METHODS:
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := formatName(method.Name)
		if method.PkgPath != "" { // method must be exported
			continue
		}

		var h callback
		h.rcvr = rcvr
		h.method = method
		h.errPos = -1

		firstArg := 1
		numIn := mtype.NumIn()
		if numIn >= 2 && mtype.In(1) == contextType {
			h.hasCtx = true
			firstArg = 2
		}

		// determine method arguments, ignore first arg since it's the receiver type
		// Arguments must be exported or builtin types
		h.argTypes = make([]reflect.Type, numIn-firstArg)
		for i := firstArg; i < numIn; i++ {
			argType := mtype.In(i)
			if !isExportedOrBuiltinType(argType) {
				continue METHODS
			}
			h.argTypes[i-firstArg] = argType
		}

		// check that all returned values are exported or builtin types
		for i := 0; i < mtype.NumOut(); i++ {
			if !isExportedOrBuiltinType(mtype.Out(i)) {
				continue METHODS
			}
		}

		// when a method returns an error it must be the last returned value
		h.errPos = -1
		for i := 0; i < mtype.NumOut(); i++ {
			if isErrorType(mtype.Out(i)) {
				h.errPos = i
				break
			}
		}

		if h.errPos >= 0 && h.errPos != mtype.NumOut()-1 {
			continue METHODS
		}

		switch mtype.NumOut() {
		case 0, 1, 2:
			if mtype.NumOut() == 2 && h.errPos == -1 { // method must one return value and 1 error
				continue METHODS
			}
			callbacks[mname] = &h
		}
	}

	return callbacks
}

var contextType = reflect.TypeOf((*context.Context)(nil)).Elem()

// 判断t是否为context类型
func isContextType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t == contextType
}

func (r *RPCX) parseArguments(argTypes []reflect.Type, args ...interface{}) ([]reflect.Value, error) {
	if len(argTypes) != len(args) {
		return nil, fmt.Errorf("args diff params")
	}
	passedArguments := make([]reflect.Value, len(args))
	var err error
	for i, v := range args {
		argTyp := argTypes[i]

		if v == nil {
			passedArguments[i] = reflect.New(argTyp).Elem()
		} else {
			passedArguments[i], err = r.parseTyp(argTyp, v, false)
			if err != nil {
				return nil, err
			}
		}
	}

	return passedArguments, nil
}

func invoke(any interface{}, name string, args ...interface{}) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn() // 输入
	if numIn > len(args) {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have minimum %d params. Have %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("Method %s must have %d params. Have %d", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("Method %s. Param[%d] must be %s. Have %s", name, i, inType, argType)
		}
	}
	return method.Call(in)[0], nil
}

func (r *RPCX) parseTyp(argTyp reflect.Type, v interface{}, isPrt bool) (param reflect.Value, err error) {
	argTypStr := argTyp.String()
	if Debug {
		logger.Info("parseType", "argTypString", argTypStr, "argTypName", argTyp.Name(), "argTypKind", argTyp.Kind())
	}
	argValue := reflect.ValueOf(v)
	if !argValue.IsValid() {
		return reflect.ValueOf(nil), fmt.Errorf("argValue is unvalid:must be %s. Have %s", argTyp, argValue.String())
	}
	paramType := argValue.Type()
	// 判断能否进行转换，如果可以转换，那么直接转换
	if paramType.ConvertibleTo(argTyp) {
		return argValue.Convert(argTyp), nil
	}

	switch val := v.(type) {
	case string:
		param, err := convutil.ToInterface(argTyp.Kind().String(), val)
		if err == nil {
			paramVal := reflect.ValueOf(param)
			if paramVal.Type().ConvertibleTo(argTyp) {
				return paramVal.Convert(argTyp), nil
			}
		}
		switch argTyp.Kind() {
		case reflect.Ptr:
			return r.parseTyp(argTyp.Elem(), v, true)
		case reflect.Struct:
			value := reflect.New(argTyp).Interface()
			var valBytes = []byte(val)
			if err := json.Unmarshal(valBytes, &value); err != nil {
				return reflect.Value{}, err
			}
			if isPrt {
				return reflect.ValueOf(value), nil
			} else {
				return reflect.ValueOf(value).Elem(), nil
			}
		default:
			return reflect.ValueOf(v), nil
		}
	default:
		return reflect.ValueOf(v), nil
	}
}

// Exec 执行
func (r *RPCX) Exec(ctx context.Context, namespace string, method string, reqId interface{}, params ...interface{}) *JsonResponse {
	namespace = formatName(namespace)
	svc1, ok := r.services.Load(namespace)
	if !ok {
		return &JsonResponse{
			Version: JsonRpcVersion,
			Id:      reqId,
			Error:   fmt.Errorf("namespace:%s no register", namespace),
		}
	}
	// 启动后，只保持读处理，所以无需锁。
	// 如果涉及到动态处理，那么需要加锁
	method = formatName(method)
	svc := svc1.(*service)
	callb1, ok := svc.callbacks.Load(method)
	if !ok {
		return &JsonResponse{
			Version: JsonRpcVersion,
			Id:      reqId,
			Error:   fmt.Errorf("namespace:%s has no [%s] method", namespace, method),
		}
	}
	callb := callb1.(*callback)
	request := &serverRequest{id: reqId, svcname: svc.name, callb: callb}
	if len(params) > 0 && len(callb.argTypes) > 0 {
		request.args, request.err = r.parseArguments(callb.argTypes, params...)
	}
	return r.handle(context.Background(), request)
}

func (r *RPCX) handle(ctx context.Context, req *serverRequest) (result *JsonResponse) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logger.Error("handle recover", "svcname", req.svcname, "err", err1)
			debug.PrintStack()
			result = &JsonResponse{
				Version: JsonRpcVersion,
				Id:      req.id,
				Error:   fmt.Errorf("%v", err1),
			}
		}
	}()
	if req.err != nil {
		return &JsonResponse{
			Version: JsonRpcVersion,
			Id:      req.id,
			Error:   req.err,
		}
	}

	// regular RPC call, prepare arguments
	if len(req.args) != len(req.callb.argTypes) {
		rpcErr := fmt.Errorf("%s%s%s expects %d parameters, got %d",
			req.svcname, ServiceMethodSeparator, req.callb.method.Name,
			len(req.callb.argTypes), len(req.args))
		return &JsonResponse{
			Version: JsonRpcVersion,
			Id:      req.id,
			Error:   rpcErr,
		}
	}

	arguments := []reflect.Value{req.callb.rcvr}
	if req.callb.hasCtx {
		arguments = append(arguments, reflect.ValueOf(ctx))
	}
	if len(req.args) > 0 {
		arguments = append(arguments, req.args...)
	}

	// [Note] execute RPC method and return result
	reply := req.callb.method.Func.Call(arguments)
	if len(reply) == 0 {
		return &JsonResponse{
			Version: JsonRpcVersion,
			Id:      req.id,
		}
	}
	if req.callb.errPos >= 0 { // test if method returned an error
		if !reply[req.callb.errPos].IsNil() {
			e := reply[req.callb.errPos].Interface().(error)
			res := &JsonResponse{
				Version: JsonRpcVersion,
				Id:      req.id,
				Error:   e,
			}
			return res
		}
	}
	return &JsonResponse{
		Version: JsonRpcVersion,
		Id:      req.id,
		Result:  reply[0].Interface(),
	}
}

func (r *RPCX) Namespaces(ctx context.Context) []*Namespace {
	if len(r.namespaces) > 0 {
		return r.namespaces
	}
	r.namespaces = make([]*Namespace, 0)
	r.services.Range(func(key, value any) bool {
		ns := key.(string)
		service := value.(*service)
		namespace := &Namespace{
			Namespace: ns,
			Methods:   make([]*Method, 0),
		}
		service.callbacks.Range(func(key, value any) bool {
			m := key.(string)
			c := value.(*callback)
			argsType := make([]string, 0, len(c.argTypes))
			for _, argType := range c.argTypes {
				argsType = append(argsType, argType.String())
			}
			namespace.Methods = append(namespace.Methods, &Method{
				Method:   m,
				ArgsType: argsType,
			})
			return true
		})
		r.namespaces = append(r.namespaces, namespace)
		return true
	})
	return r.namespaces
}

func (r *RPCX) Namespaces2(ctx context.Context) map[string]map[string]*Method {
	if len(r.namespacesMap) > 0 {
		return r.namespacesMap
	}
	r.namespacesMap = make(map[string]map[string]*Method, 0)
	r.services.Range(func(key, value any) bool {
		ns := key.(string)
		service := value.(*service)
		methods := make(map[string]*Method, 0)
		service.callbacks.Range(func(key, value any) bool {
			m := key.(string)
			c := value.(*callback)
			argsType := make([]string, 0, len(c.argTypes))
			for _, argType := range c.argTypes {
				argsType = append(argsType, argType.String())
			}
			methods[m] = &Method{
				Method:   m,
				ArgsType: argsType,
			}
			return true
		})
		r.namespacesMap[ns] = methods
		return true
	})
	return r.namespacesMap
}
