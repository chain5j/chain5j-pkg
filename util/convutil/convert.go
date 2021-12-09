// Package convutil
//
// @author: xwc1125
package convutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func BytesToBool(input []byte) (bool, error) {
	if input == nil || len(input) == 0 {
		return false, nil
	}
	return input[0] == 1, nil
}

func BytesToString(input []byte) (string, error) {
	if input == nil || len(input) == 0 {
		return "", nil
	}
	return string(input), nil
}

func BytesToInt64(input []byte) (int64, error) {
	if input == nil || len(input) == 0 {
		return 0, nil
	}
	if len(input) == 3 {
		input = append([]byte{0}, input...)
	}
	bytesBuffer := bytes.NewBuffer(input)
	switch len(input) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int64(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int64(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int64(tmp), err
	case 8:
		var tmp int64
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return tmp, err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt64 bytes lenth is invaild!")
	}
}

func BytesToUint64(input []byte) (uint64, error) {
	if input == nil || len(input) == 0 {
		return 0, nil
	}
	if len(input) == 3 {
		input = append([]byte{0}, input...)
	}
	bytesBuffer := bytes.NewBuffer(input)
	switch len(input) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return uint64(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return uint64(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return uint64(tmp), err
	case 8:
		var tmp uint64
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return tmp, err
	default:
		return 0, fmt.Errorf("%s", "BytesToUint64 bytes lenth is invaild!")
	}
}
