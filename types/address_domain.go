// Package types
//
// @author: xwc1125
package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/chain5j/chain5j-pkg/codec/rlp"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
)

var (
	EmptyDomainAddress = DomainAddress{
		Addr: EmptyAddress,
	}
)

type DomainAddress struct {
	Addr       Address `json:"addr"`       // address
	DomainAddr string  `json:"domainAddr"` // domainAddr
}

func NewDomainAddress(domain string, addr Address) DomainAddress {
	return DomainAddress{
		Addr:       addr,
		DomainAddr: domain,
	}
}

// DomainToAddress domainAddr转Address。addr格式：xxx@chain5j.com或0xxxxxxxx
func DomainToAddress(s string) Address {
	if strings.Contains(s, "@") {
		split := strings.Split(s, "@")
		if hexutil.IsHex(split[0]) {
			return BytesToAddress(hexutil.MustDecode(split[0]))
		}
		return StringToAddress(split[0])
	}
	return BytesToAddress(hexutil.MustDecode(s))
}

// FromDomainAddress domainAddress 转换为DomainAddress对象。addr格式：xxx@chain5j.com:0xxxxxxxx
func FromDomainAddress(addr string) DomainAddress {
	if strings.Contains(addr, ":") {
		split := strings.Split(addr, ":")
		if len(split) != 2 && len(split) != 2 {
			return EmptyDomainAddress
		}
		var a DomainAddress
		a.DomainAddr = split[0]
		if len(split) == 2 {
			a.Addr = BytesToAddress(hexutil.MustDecode(split[1]))
		}
		return a
	}
	return EmptyDomainAddress
}

func (a DomainAddress) Len() int {
	return len(a.DomainAddr)
}

func (a DomainAddress) Bytes() []byte {
	return []byte(a.String())
}

func (a DomainAddress) FromBytes(b []byte) (Addr, error) {
	addr := string(b)
	return a.FromStr(addr)
}

func (a DomainAddress) String() string {
	return a.DomainAddr + ":" + a.Addr.Hex()
}

func (a DomainAddress) FromStr(addr string) (Addr, error) {
	if strings.Contains(addr, ":") {
		split := strings.Split(addr, ":")
		if len(split) != 2 && len(split) != 2 {
			return EmptyDomainAddress, errors.New("invalid bytes")
		}
		a.DomainAddr = split[0]
		if len(split) == 2 {
			a.Addr = BytesToAddress(hexutil.MustDecode(split[1]))
		}
		return a, nil
	}
	return EmptyDomainAddress, errors.New("invalid bytes")
}

func (a DomainAddress) Nil() bool {
	if a == EmptyDomainAddress {
		return true
	}
	if a.DomainAddr == "" && a.Addr == EmptyAddress {
		return true
	}
	return false
}

func (a DomainAddress) Validate(addr string) bool {
	return !a.Nil()
}

func (a DomainAddress) Hash() Hash {
	return BytesToHash(a.Bytes())
}

func (a DomainAddress) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a.Bytes())
}

func (a DomainAddress) Big() *big.Int {
	return new(big.Int).SetBytes(a.Bytes())
}

func (a *DomainAddress) UnmarshalJSON(data []byte) error {
	type erased DomainAddress
	e := erased{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		return err
	}

	a.Addr = e.Addr
	a.DomainAddr = e.DomainAddr
	return nil
}

func (a DomainAddress) Value() (driver.Value, error) {
	return a.Bytes(), nil
}

func (a *DomainAddress) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, a.Bytes())
}

func (a *DomainAddress) DecodeRLP(s *rlp.Stream) error {
	var bytes []byte
	err := s.Decode(&bytes)
	if err != nil {
		return err
	}
	addr := string(bytes)
	domainAddress := FromDomainAddress(addr)
	a.Addr = domainAddress.Addr
	a.DomainAddr = domainAddress.DomainAddr
	return nil
}
