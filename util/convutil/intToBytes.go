// Package convutil
//
// @author: xwc1125
// @date: 2020/2/15
package convutil

import "encoding/binary"

// Int64ToBytes int64 转 byte
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BytesToInt64 byte 转 int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// Uint64ToBytes int64 转 byte
func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

// BytesToUint64 byte 转 uint64
func BytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}
