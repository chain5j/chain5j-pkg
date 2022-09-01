// Package icap
//
// @author: xwc1125
package icap

import (
	"github.com/chain5j/chain5j-pkg/codec/json"
	"github.com/chain5j/chain5j-pkg/crypto/base/base36"
	"github.com/chain5j/chain5j-pkg/types"
)

// THash 交易Hash
type THash struct {
	ChainName string     `json:"chainName" mapstructure:"chainName"` // 链名称
	Index     uint64     `json:"index" mapstructure:"index"`         // dbindex
	Hash      types.Hash `json:"hash" mapstructure:"hash"`           // 交易hash值
}

// THashFromIban 字符串转对象
func THashFromIban(tHash string) (THash, error) {
	iBanInfo := NewIBanInfo(4, 0, 64, tHash)
	customer, err := ParseICAP(*iBanInfo)
	if err != nil {
		return THash{}, err
	}
	hash := types.BytesToHash(customer.Customer())
	return THash{
		ChainName: customer.Currency(),
		Index:     0,
		Hash:      hash,
	}, nil
}

// String tHash的值
func (h THash) String() string {
	indexBase36 := base36.Encode(uint64(h.Index))
	customer := NewCustomer(
		h.ChainName,
		LeftJoin(indexBase36, 4),
		60,
		h.Hash.Bytes())
	iBanInfo, _ := ToICAP(*customer)
	return iBanInfo.Iban()
}

func (h THash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h THash) Bytes() []byte {
	bytes, _ := json.Marshal(h.String())
	return bytes
}

func (h THash) Nil() bool {
	if h.ChainName == "" || h.Hash.Nil() {
		return true
	}
	return false
}
