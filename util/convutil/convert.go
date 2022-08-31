// Package convutil
//
// @author: xwc1125
package convutil

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
