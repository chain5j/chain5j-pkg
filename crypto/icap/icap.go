// Package icap
//
// @author: xwc1125
// @date: 2021/4/14
package icap

import (
	"errors"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/chain5j/chain5j-pkg/crypto/base/base36"
	"github.com/chain5j/chain5j-pkg/types"
)

var (
	errICAPLength      = errors.New("invalid ICAP length")
	errICAPEncoding    = errors.New("invalid ICAP encoding")
	errICAPChecksum    = errors.New("invalid ICAP checksum")
	errICAPCountryCode = errors.New("invalid ICAP country code")
	errICAPAssetIdent  = errors.New("invalid ICAP asset identifier")
	errICAPInstCode    = errors.New("invalid ICAP institution code")
	errICAPClientIdent = errors.New("invalid ICAP client identifier")
)

var (
	Big1  = big.NewInt(1)
	Big97 = big.NewInt(97)
	Big98 = big.NewInt(98)
)

func ToICAP(customer Customer) (*IBanInfo, error) {
	enc := base36.EncodeBytes(customer.Customer())
	currencyLen := len(customer.Currency())
	orgCodeLen := len(customer.OrgCode())
	// 减去2位校验码
	interLen := customer.ResultLen() - currencyLen - orgCodeLen - 2 - len(enc)
	if interLen > 0 {
		enc = join(strings.Repeat("0", interLen), enc)
	}
	icap := join(customer.Currency(), checkDigits(customer.Currency(), customer.OrgCode(), enc), customer.OrgCode(), enc)
	return NewIBanInfo(currencyLen, orgCodeLen, len(customer.customer), icap), nil
}

func ParseICAP(iban IBanInfo) (*Customer, error) {
	if err := ValidCheckSumWithLen(iban.currencyLen, iban.orgCodeLen, iban.iban); err != nil {
		return nil, err
	}
	// checksum is ISO13616, Ethereum address is base-36
	l := iban.currencyLen + 2 + iban.orgCodeLen
	// bigAddr, _ := new(big.Int).SetString(iban.iban[l:], 36)
	// hex := hexutil.EncodeBig(bigAddr)
	bytes := base36.DecodeToBytes(true, iban.iban[l:])
	return &Customer{
		currency:  iban.iban[:iban.currencyLen],
		orgCode:   iban.iban[iban.currencyLen+2 : iban.currencyLen+2+iban.orgCodeLen],
		resultLen: len(bytes) + 2,
		customer:  bytes,
	}, nil
}

//export ConvertAddressToICAP
func ConvertAddressToICAP(prefix string, orgCode string, a types.Address) string {
	prefix = strings.ToUpper(prefix)
	enc := base36.EncodeBytes(a.Bytes())
	// zero padd encoded address to Direct ICAP length if needed
	if len(enc) < 31 {
		enc = join(strings.Repeat("0", 31-len(enc)), enc)
	}
	icap := join(prefix, checkDigits(prefix, orgCode, enc), orgCode, enc)

	l := 31 + len(orgCode) + len(prefix) + 2
	if len(icap) != l {
		log.Println("生成的地址出错", "addr", a.Hex())
	}
	return strings.ToLower(icap)
}

//export ConvertICAPToAddress
func ConvertICAPToAddress(prefix string, orgCodeLen int, s string) (types.Address, error) {
	prefix = strings.ToUpper(prefix)
	s = strings.ToUpper(s)
	l := 31 + orgCodeLen + len(prefix) + 2
	switch len(s) {
	case 35: // "XE" + 2 digit checksum + 31 base-36 chars of address
		return parseICAP(prefix, s)
	case 34: // "XE" + 2 digit checksum + 30 base-36 chars of address
		return parseICAP(prefix, s)
	case 20: // "XE" + 2 digit checksum + 3-char asset identifier +
		// 4-char institution identifier + 9-char institution client identifier
		return parseIndirectICAP(prefix, s)
	case l: // "prefix" + 2 digit checksum + orgCodeLen + 3-char asset identifier +
		// 4-char institution identifier + 9-char institution client identifier
		return parseSelfICAP(prefix, orgCodeLen, s)
	default:
		return types.Address{}, errICAPLength
	}
}
func parseSelfICAP(prefix string, orgCodeLen int, s string) (types.Address, error) {
	if !strings.HasPrefix(s, prefix) {
		return types.Address{}, errICAPCountryCode
	}
	if orgCodeLen > 0 {
		if err := ValidCheckSumWithLen(len(prefix), orgCodeLen, s); err != nil {
			return types.Address{}, err
		}
	} else {
		if err := validCheckSum(s); err != nil {
			return types.Address{}, err
		}
	}
	// checksum is ISO13616, Ethereum address is base-36
	l := len(prefix) + 2 + orgCodeLen
	bigAddr, _ := new(big.Int).SetString(s[l:], 36)
	return types.BigToAddress(bigAddr), nil
}
func parseICAP(prefix string, s string) (types.Address, error) {
	if !strings.HasPrefix(s, prefix) {
		return types.Address{}, errICAPCountryCode
	}
	if err := validCheckSum(s); err != nil {
		return types.Address{}, err
	}
	// checksum is ISO13616, Ethereum address is base-36
	bigAddr, _ := new(big.Int).SetString(s[4:], 36)
	return types.BigToAddress(bigAddr), nil
}
func parseIndirectICAP(prefix string, s string) (types.Address, error) {
	if !strings.HasPrefix(s, prefix) {
		return types.Address{}, errICAPCountryCode
	}
	if s[4:7] != "ETH" {
		return types.Address{}, errICAPAssetIdent
	}
	if err := validCheckSum(s); err != nil {
		return types.Address{}, err
	}
	return types.Address{}, errors.New("not implemented")
}

// https://en.wikipedia.org/wiki/International_Bank_Account_Number#Validating_the_IBAN
func validCheckSum(s string) error {
	s = join(s[4:], s[:4])
	expanded, err := iso13616Expand(s)
	if err != nil {
		return err
	}
	checkSumNum, _ := new(big.Int).SetString(expanded, 10)
	if checkSumNum.Mod(checkSumNum, Big97).Cmp(Big1) != 0 {
		return errICAPChecksum
	}
	return nil
}

func ValidCheckSumWithLen(prefixLen, orgCodeLen int, s string) error {
	// s=prefix+check+orgCode+addr
	l := prefixLen + 2
	s1 := join(s[l:], s[:prefixLen], s[prefixLen:prefixLen+2]) // orgCode+addr+prefix+check
	expanded, err := iso13616Expand(s1)
	if err != nil {
		return err
	}
	checkSumNum, _ := new(big.Int).SetString(expanded, 10)
	if checkSumNum.Mod(checkSumNum, Big97).Cmp(Big1) != 0 {
		return errICAPChecksum
	}
	return nil
}

func checkDigits(prefix, orgCode, s string) string {
	prefix = strings.ToUpper(prefix)
	expanded, _ := iso13616Expand(join(orgCode, s, prefix+"00")) // orgCode+addr+prefix+00
	num, _ := new(big.Int).SetString(expanded, 10)
	num.Sub(Big98, num.Mod(num, Big97))

	checkDigits := num.String()
	// zero padd checksum
	if len(checkDigits) == 1 {
		checkDigits = join("0", checkDigits)
	}
	return checkDigits
}

// not base-36, but expansion to decimal literal: A = 10, B = 11, ... Z = 35
func iso13616Expand(s string) (string, error) {
	var parts []string
	for _, c := range s {
		i := uint64(c)
		// 0-9 or A-Z
		if i < 48 || (i > 57 && i < 65) || i > 90 {
			return "", errICAPEncoding
		}

		if i > 97 {
			parts = append(parts, strconv.FormatUint(uint64(c)-87, 10))
		} else if i >= 65 {
			parts = append(parts, strconv.FormatUint(uint64(c)-55, 10))
		} else {
			parts = append(parts, string(c))
		}
	}
	return join(parts...), nil
}

func join(s ...string) string {
	return strings.Join(s, "")
}

func LeftJoin(str string, resultLen int) string {
	if len(str) < resultLen {
		str = join(strings.Repeat("0", resultLen-len(str)), str)
	}
	return str
}
