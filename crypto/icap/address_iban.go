// Package icap
//
// @author: xwc1125
package icap

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/types"
)

var (
	AddressPrettyPrefix = 4
	AddressPrettyLen    = 37
	EmptyAddressPretty  = addressPretty{
		Addr: types.EmptyAddress,
	}
)

// addressPretty 优雅地址，base36类型
type addressPretty struct {
	ChainName string        `json:"chainName" mapstructure:"chainName"` // 链名称
	Addr      types.Address `json:"addr" mapstructure:"addr"`           // 地址
	ibanAddr  string        // ibanAddr地址
}

func NewAddressPretty(chainName string, address types.Address) addressPretty {
	customer := NewCustomer(
		chainName,
		"",
		AddressPrettyLen,
		address.Hex())
	iBanInfo, err := ToICAP(*customer)
	if err != nil {
		return EmptyAddressPretty
	}
	iban := iBanInfo.Iban()

	return addressPretty{
		ChainName: chainName,
		Addr:      address,
		ibanAddr:  iban,
	}
}

// AddressPrettyFromIban 字符串转对象
func AddressPrettyFromIban(tAddr string) (addressPretty, error) {
	iBanInfo := NewIBanInfo(AddressPrettyPrefix, 0, 31, tAddr)
	customer, err := ParseICAP(*iBanInfo)
	if err != nil {
		return addressPretty{}, err
	}
	address := types.HexToAddress(customer.Customer())
	return addressPretty{
		ChainName: customer.Currency(),
		Addr:      address,
		ibanAddr:  tAddr,
	}, nil
}

func (a addressPretty) Len() int { return AddressPrettyLen }

func (a addressPretty) Bytes() []byte { return []byte(a.String()) }

func (a addressPretty) FromBytes(b []byte) (types.Addr, error) {
	addr := string(b)
	return a.FromStr(addr)
}

func (a addressPretty) String() string {
	return a.ibanAddr
}

func (a addressPretty) FromStr(addr string) (types.Addr, error) {
	return AddressPrettyFromIban(addr)
}

func (a addressPretty) Validate(addr string) bool {
	_, err := AddressPrettyFromIban(addr)
	if err != nil {
		return false
	}
	return true
}

func (a addressPretty) Nil() bool {
	return a.Addr == types.EmptyAddress
}

func (a addressPretty) Hash() types.Hash {
	return types.BytesToHash(a.Bytes())
}

func (a addressPretty) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a.Bytes())
}

// TerminalString String implements fmt.Stringer.
func (a addressPretty) TerminalString() string {
	return a.String()
}

func (a addressPretty) MarshalJSON() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *addressPretty) UnmarshalJSON(input []byte) error {
	addr := string(input)
	pretty, err := AddressPrettyFromIban(addr)
	if err != nil {
		return err
	}
	a.ChainName = pretty.ChainName
	a.Addr = pretty.Addr
	a.ibanAddr = pretty.ibanAddr
	return nil
}
