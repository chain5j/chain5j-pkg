// Package types
//
// @author: xwc1125
package types

const TxTypeUnknown = "UNKNOWN"

// TxType the transaction type
type TxType string

// Value value of txType
func (txType *TxType) Value() string {
	return string(*txType)
}

// ValueOf string to txType
func (txType *TxType) ValueOf(v string) TxType {
	return TxType(v)
}

// TxTypes the array of txType
type TxTypes []TxType

// Len len
func (types TxTypes) Len() int {
	return len(types)
}

// Less less
func (types TxTypes) Less(i, j int) bool {
	if types[i] < types[j] {
		return true
	} else {
		return false
	}
}

// Swap swap
func (types TxTypes) Swap(i, j int) {
	types[i], types[j] = types[j], types[i]
}
