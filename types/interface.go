// Package types
// 
// @author: xwc1125
// @date: 2021/6/21
package types

import "fmt"

type Addr interface {
	Len() int                          // 地址长度
	Bytes() []byte                     // 地址bytes
	FromBytes(b []byte) (Addr, error)  // 将bytes转换为address对象
	String() string                    // 地址转换为string [implements fmt.Stringer]
	FromStr(addr string) (Addr, error) // 将字符串转换为address对象

	Validate(addr string) bool // 验证地址正确性
	Nil() bool                 // 判断地址是否为空
	Hash() Hash                // 地址的hash

	Format(s fmt.State, c rune) // 地址格式化
}
