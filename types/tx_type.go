// description: chain5j-core 
// 
// @author: xwc1125
// @date: 2020/10/16
package types

// 交易类型
type TxType string

func (txType *TxType) Value() string {
	return string(*txType)
}

func (txType *TxType) ValueOf(v string) TxType {
	return TxType(v)
}

type TxTypes []TxType

func (types TxTypes) Len() int {
	return len(types)
}

func (types TxTypes) Less(i, j int) bool {
	if types[i] < types[j] {
		return true
	} else {
		return false
	}
}

func (types TxTypes) Swap(i, j int) {
	types[i], types[j] = types[j], types[i]
}
