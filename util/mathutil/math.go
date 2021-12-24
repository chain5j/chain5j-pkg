// Package mathutil
//
// @author: xwc1125
package mathutil

// If 三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
