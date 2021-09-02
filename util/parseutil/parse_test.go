// Package parseutil
// 
// @author: xwc1125
// @date: 2020/11/3
package parseutil

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	//调用无参方法
	values := ReflectInterface(testReflectParam, nil)
	for i, v := range values {
		fmt.Println(i, v)
	}
	//调用有参方法
	values = ReflectInterface(testReflectParam2, GetValues(5, "Hello"))
	for i, v := range values {
		fmt.Println(i, v)
	}
}

//无参方法
func testReflectParam() (string, string) {
	return "hello world", "你好！"
}

//有参方法
func testReflectParam2(i int, s string) (int, string) {
	i++
	return i, s
}
