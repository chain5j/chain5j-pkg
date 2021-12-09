// Package gmsm
//
// @author: xwc1125
package gmsm

import (
	"crypto/elliptic"
	"github.com/tjfoc/gmsm/sm2"
)

// Curve 椭圆曲线
func Curve() elliptic.Curve {
	return sm2.P256Sm2()
}
