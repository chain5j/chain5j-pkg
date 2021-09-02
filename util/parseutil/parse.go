// Package parseutil
//
// @author: xwc1125
// @date: 2019/10/29
package parseutil

import (
	"log"
	"reflect"
)

func GetValue(param interface{}) reflect.Value {
	return reflect.ValueOf(param)
}

//GetValues 根据参数获取对应的Values
func GetValues(param ...interface{}) []reflect.Value {
	vals := make([]reflect.Value, 0, len(param))
	for i := range param {
		vals = append(vals, reflect.ValueOf(param[i]))
	}
	return vals
}

func ReflectInterface(funcInter interface{}, paramsValue []reflect.Value) []reflect.Value {
	v := reflect.ValueOf(funcInter)
	if v.Kind() != reflect.Func {
		log.Fatal("funcInter is not func")
	}

	values := v.Call(paramsValue) //方法调用并返回值
	return values
}
