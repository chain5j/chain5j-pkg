// description: chain5j-core 
// 
// @author: xwc1125
// @date: 2020/10/17
package types

type TxUuid []byte

func (txUuid *TxUuid) Value() []byte {
	return ([]byte)(*txUuid)
}

func (txUuid *TxUuid) ValueOf(v []byte) TxUuid {
	return TxUuid(v)
}

