// Package base64 base64编码
//
// @author: xwc1125
package base64

import "encoding/base64"

// Encode base64编码
func Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// Decode base64解码
func Decode(enc string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(enc)
}
