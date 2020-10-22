// description: chain5j 
// 
// @author: xwc1125
// @date: 2020/10/19
package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/chain5j/chain5j-pkg/codec/rlp"
	log "github.com/chain5j/log15"
	"io"
	"math/big"
	"strings"
)

var (
	EmptyDomainAddress = DomainAddress{}
	emptyAddress       = Address{}
)

type DomainAddress struct {
	Addr   Address `json:"addr"`   // 地址
	Domain string  `json:"domain"` // 域
}

func AddressToDomainAddress(addr Address) DomainAddress {
	// TODO 可以通过address获取域
	return DomainAddress{
		Addr:   addr,
		Domain: addr.Hex(),
	}
}

func BytesToDomainAddress(addr []byte) DomainAddress {
	address := BytesToAddress(addr)
	return AddressToDomainAddress(address)
}

func BigToDomainAddress(b *big.Int) DomainAddress {
	address := BigToAddress(b)
	return AddressToDomainAddress(address)
}

func HexToDomainAddress(s string) DomainAddress {
	address := HexToAddress(s)
	if address.Nil() {
		return EmptyDomainAddress
	}
	return AddressToDomainAddress(address)
}

func DomainToDomainAddress(domain string) DomainAddress {
	return DomainAddress{
		Domain: domain,
	}
}

func (addr DomainAddress) IsEmpty() bool {
	if addr == EmptyDomainAddress {
		return true
	}
	if addr.Domain == "" && addr.Addr == emptyAddress {
		return true
	}
	return false
}
func (addr DomainAddress) Bytes() []byte {
	if addr == EmptyDomainAddress {
		return nil
	}
	if addr.Addr != emptyAddress {
		return addr.Addr[:]
	}
	return []byte(addr.Domain)
}

func (a DomainAddress) Big() *big.Int {
	return new(big.Int).SetBytes(a.Bytes())
}

func (a *DomainAddress) UnmarshalJSON(data []byte) error {
	type erased DomainAddress
	e := erased{}
	err := json.Unmarshal(data, &e)
	if err == nil {
		a.Addr = e.Addr
		a.Domain = e.Domain
		return nil
	}
	var input string
	err = json.Unmarshal(data, &input)
	if err != nil {
		return err
	}
	if strings.HasPrefix(input, "0x") && len(input) == 42 {
		addr := HexToDomainAddress(input)
		a.Addr = addr.Addr
		a.Domain = addr.Domain
		return nil
	} else if !strings.HasPrefix(input, "0x") && len(input) == 40 {
		addr := HexToDomainAddress(input)
		a.Addr = addr.Addr
		a.Domain = addr.Domain
		return nil
	} else {
		addr := DomainToDomainAddress(input)
		a.Addr = addr.Addr
		a.Domain = addr.Domain
		return nil
	}
}

func (a DomainAddress) String() string {
	return a.Domain
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
		log.Error("DecodeTxsRLP", "err", err)
		return err
	}
	domainAddress := BytesToDomainAddress(bytes)
	a = &domainAddress
	return nil
}
