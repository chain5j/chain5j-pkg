// Package reflectutil
//
// @author: xwc1125
package reflectutil

import (
	"github.com/chain5j/logger"
	"log"
	"reflect"
)

// ToPointer 将 interface{T} 转变为 interface{*T}
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

// GetFieldName structName 的 type 不是结构体类型，就会报以下错误：panic: reflect: NumField of non-struct type，故需要在程序中加以判断。
//获取结构体中字段的名称
func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

func GetValueByFieldName(structName interface{}, fieldName string) interface{} {
	t := reflect.ValueOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		logger.Println("Check type error not Struct")
		return nil
	}
	fieldByName := t.FieldByName(fieldName)
	return fieldByName.Interface()
}

//获取结构体中Tag的值，如果没有tag则返回字段值
func GetTagName(structName interface{}, tagKey string) []FieldInfo {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	fieldInfos := make([]FieldInfo, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		fieldName := t.Field(i).Name
		fieldType := t.Field(i).Type.Name() // 类型

		fieldInfo := FieldInfo{
			FieldName: fieldName,
			FieldType: fieldType,
		}
		// 获取tag内容
		tagValue := t.Field(i).Tag.Get(tagKey)
		fieldInfo.TagValue = tagValue
		fieldInfos = append(fieldInfos, fieldInfo)
	}
	return fieldInfos
}

type FieldInfo struct {
	FieldName string // 字段名称
	FieldType string // 字段类型
	TagValue  string // tag的值
}
