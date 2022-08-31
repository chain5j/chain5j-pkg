// Package copier
//
// @author: xwc1125
package copier

import (
	"github.com/jinzhu/copier"
)

// tag:copier:"must,nopanic"
// tag:copier:"-"
// tag:copier:"EmployeNum"// change field name

const (
	String  = copier.String
	Bool    = copier.Bool
	Int     = copier.Int
	Float32 = copier.Float32
	Float64 = copier.Float64
)

type (
	Option        copier.Option
	TypeConverter copier.TypeConverter
)

var (
	Copy           = copier.Copy
	CopyWithOption = copier.CopyWithOption
)
