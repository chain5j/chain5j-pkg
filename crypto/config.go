// description: Atom 
// 
// @author: xwc1125
// @date: 2019/9/9
package crypto

import (
	"crypto/elliptic"
	"github.com/chain5j/chain5j-pkg/crypto/secp256k1"
	"github.com/tjfoc/gmsm/sm2"
)

const (
	P256    = "P-256"
	S256    = "S-256"
	SM2P256 = "SM2-P-256"
	P384    = "P-384"
	P521    = "P-521"
)

var CryptoName = P256

func DefaultCryptoType() elliptic.Curve {
	return elliptic.P256()
}

func CurveType(curveType string) elliptic.Curve {
	switch curveType {
	case P256:
		return elliptic.P256()
	case S256:
		return secp256k1.S256()
	case SM2P256:
		return sm2.P256Sm2()
	}
	return elliptic.P256()
}

func CurveName(curve elliptic.Curve) string {
	if curve == secp256k1.S256() {
		return S256
	}
	name := curve.Params().Name
	return name
}
