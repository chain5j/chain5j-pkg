// Package icap
// 
// @author: xwc1125
// @date: 2021/6/2
package icap

import "strings"

type Customer struct {
	currency  string `json:"currency"`  // 资产标识符
	orgCode   string `json:"orgCode"`   // 机构标识符
	resultLen int    `json:"resultLen"` // 机构客户标识符base36的标准长度+currencyLen+orgCodeLen+2位校验码 如果长度不足时，base36会进行补0
	customer  string `json:"customer"`  // 机构客户标识符[16进制]
}

type IBanInfo struct {
	currencyLen int    `json:"currencyLen"` // 资产标识符长度
	orgCodeLen  int    `json:"orgCodeLen"`  // 机构标识符长度
	customerLen int    `json:"customerLen"` // 机构客户标识符的标准长度
	iban        string `json:"iban"`        // iban标识
}

func NewCustomer(currency, orgCode string, resultLen int, customer string) *Customer {
	return &Customer{
		currency:  currency,
		orgCode:   orgCode,
		resultLen: resultLen,
		customer:  customer,
	}
}
func NewIBanInfo(currencyLen, orgCodeLen, customerLen int, iban string) *IBanInfo {
	return &IBanInfo{
		currencyLen: currencyLen,
		orgCodeLen:  orgCodeLen,
		customerLen: customerLen,
		iban:        iban,
	}
}

func (c *Customer) Currency() string {
	return strings.ToUpper(c.currency)
}

func (c *Customer) OrgCode() string {
	return strings.ToUpper(c.orgCode)
}

func (c *Customer) ResultLen() int {
	return c.resultLen
}

func (c *Customer) Customer() string {
	return c.customer
}

func (b *IBanInfo) CurrencyLen() int {
	return b.currencyLen
}

func (b *IBanInfo) OrgCodeLen() int {
	return b.orgCodeLen
}

func (b *IBanInfo) CustomerLen() int {
	return b.customerLen
}

func (b *IBanInfo) Iban() string {
	return strings.ToUpper(b.iban)
}
