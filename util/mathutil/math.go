// description: chain5j 
// 
// @author: xwc1125
// @date: 2020/11/19
package mathutil

// 三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	} else {
		return falseVal
	}
}
