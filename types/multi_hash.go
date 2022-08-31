// Package types
//
// @author: xwc1125
package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

var (
	EmptyMultiHash = MultiHash{}
)

type MultiHash []byte

func BytesToMultiHash(b []byte) MultiHash {
	var h MultiHash
	h.SetBytes(b)
	return h
}

func StringToMultiHash(s string) MultiHash {
	s1, err := strconv.Unquote(s)
	if err == nil {
		s = s1
	}
	return BytesToMultiHash([]byte(s))
}

func (h MultiHash) Bytes() []byte { return h[:] }

func (h MultiHash) String() string {
	return string(h.Bytes())
}

func (h MultiHash) TerminalString() string {
	return h.String()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
// log会调用此接口
func (h MultiHash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h.Bytes())
}

func (h *MultiHash) SetBytes(b []byte) {
	var des = make([]byte, len(b))
	copy(des, b)
	*h = des
}

func (h MultiHash) Nil() bool {
	return len(h) == 0
}

type extMultiHash string

func (h MultiHash) MarshalJSON() ([]byte, error) {
	var extMultiHash = h.String()
	return json.Marshal(extMultiHash)
}
func (h *MultiHash) UnmarshalJSON(data []byte) error {
	var extMultiHash extMultiHash
	err := json.Unmarshal(data, &extMultiHash)
	if err != nil {
		return err
	}
	multiHash := StringToMultiHash(string(extMultiHash))
	*h = multiHash
	return nil
}

func (h MultiHash) MarshalText() ([]byte, error) {
	return h.MarshalJSON()
}

// UnmarshalText .
func (h *MultiHash) UnmarshalText(input []byte) error {
	return h.UnmarshalJSON(input)
}

func (h MultiHash) Big() *big.Int {
	return new(big.Int).SetBytes(h[:])
}

// func (h MultiHash) EncodeRLP(w io.Writer) error {
// 	return rlp.Encode(w, h.String())
// }
// func (h *MultiHash) DecodeRLP(s *rlp.Stream) error {
// 	origin, err := s.Raw()
// 	if err == nil {
// 		var hashStr string
// 		err = rlp.DecodeBytes(origin, &hashStr)
// 		if err != nil {
// 			return err
// 		}
// 		multiHash := StringToMultiHash(hashStr)
// 		*h = multiHash
// 	}
// 	return err
// }
