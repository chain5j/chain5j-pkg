// Package icap
//
// @author: xwc1125
package icap

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/types"
	"log"
	"strings"
	"testing"
)

func TestConvertAddressToICAP(t *testing.T) {
	icap := ConvertAddressToICAP("xwc1125", "0000", types.HexToAddress("0x0000059B6D857f25fBDc1ab0Bb78AAa048f638A5"))
	fmt.Println("icap", icap)
	address, err := ConvertICAPToAddress("xwc1125", 4, strings.ToUpper(icap))
	if err != nil {
		panic(err)
	}

	fmt.Println("address", address.Hex())

	address, err = ConvertICAPToAddress("xwc1125", 4, strings.ToLower(icap))
	if err != nil {
		panic(err)
	}

	fmt.Println("address2", address.Hex())
}

func TestIban(t *testing.T) {
	banInfo, err := ToICAP(Customer{
		currency:  "chain5j",
		orgCode:   "0000",
		resultLen: 60,
		customer:  "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	})
	if err != nil {
		panic(err)
	}
	customer, err := ParseICAP(*banInfo)
	if err != nil {
		panic(err)
	}
	log.Println(customer)
}
