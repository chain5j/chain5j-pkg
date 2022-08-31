// Package errorx
//
// @author: xwc1125
package errorx

import (
	"fmt"

	"github.com/chain5j/chain5j-pkg/codec/json"
)

type Error interface {
	Error() string  // 错误信息
	ErrorCode() int // 错误码
}

type RespInfo struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Errorf returns Error(code, fmt.Sprintf(format, a...)).
func Errorf(code int, format string, a ...interface{}) RespInfo {
	return Errorw(code, fmt.Sprintf(format, a...))
}

func Errorw(code int, msg string) RespInfo {
	return RespInfo{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
}

func Errordf(code int, data interface{}, format string, a ...interface{}) RespInfo {
	return Errord(code, data, fmt.Sprintf(format, a...))
}

func Errord(code int, data interface{}, msg string) RespInfo {
	return RespInfo{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func (err RespInfo) String() string {
	marshal, _ := json.Marshal(err)
	return string(marshal)
}

func (err RespInfo) Error() string {
	var errMsg string
	if err.Message == "" {
		errMsg = fmt.Sprintf("code=%d", err.Code)
	} else {
		errMsg = fmt.Sprintf("code=%d, msg=%s", err.Code, err.Message)
	}
	if err.Data != nil {
		errMsg = fmt.Sprintf("%s, data=%v", errMsg, err.Data)
	}
	return errMsg
}

func (err RespInfo) ErrorCode() int {
	return err.Code
}

func (err RespInfo) ErrorData() interface{} {
	return err.Data
}
