// Package icap
//
// @author: xwc1125
// @date: 2021/4/14
package icap

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/types"
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
