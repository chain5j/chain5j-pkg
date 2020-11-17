// description: xblog
//
// @author: xwc1125
// @date: 2020/3/11
package reflectutil

import (
	"reflect"
)

// 将 interface{T} 转变为 interface{*T}
// 如果是指针，直接返回
func ToPointer(w interface{}) interface{} {
	typeOf := reflect.TypeOf(w)
	if typeOf.Kind() == reflect.Ptr {
		return w
	}
	//reflect.PtrTo(typeOf)
	valueOf := reflect.ValueOf(w)
	pv := reflect.New(typeOf)
	pv.Elem().Set(valueOf)
	return pv.Interface()
}

func DelPointer(w interface{}) interface{} {
	typeOf := reflect.TypeOf(w)
	if typeOf.Kind() != reflect.Ptr {
		return w
	}
	valueOf := reflect.ValueOf(w)
	return valueOf.Elem().Interface()
}
