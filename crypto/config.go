// description: Atom 
// 
// @author: xwc1125
// @date: 2019/9/9
package crypto

import (
	"crypto/elliptic"
	"github.com/chain5j/chain5j-pkg/crypto/secp256k1"
	"github.com/tjfoc/gmsm/sm2"
	"strings"
)

const (
	P256    = "P256"
	S256    = "S256"
	SM2P256 = "SM2P256"
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

func CurveName(params *elliptic.CurveParams) string {
	name := params.Name
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "", "")
	return name
}
