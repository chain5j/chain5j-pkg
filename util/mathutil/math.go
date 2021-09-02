// Package mathutil
//
// @author: xwc1125
// @date: 2020/11/19
package mathutil

// If 三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
