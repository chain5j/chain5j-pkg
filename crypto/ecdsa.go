// Package crypto
//
// @author: xwc1125
// @date: 2021/8/19
package crypto

import (
	"crypto/ecdsa"
	"github.com/chain5j/chain5j-pkg/codec"
)

// PrivateKey 私钥接口
type PrivateKey interface {
	// GenerateECDSAKey 生成ecdsa key
	GenerateECDSAKey() (*ecdsa.PrivateKey, error)
	// GenerateKey 生成crypto
	GenerateKey() (PrivateKey, error)
	// ToECDSA 转成ecdsa
	ToECDSA() *ecdsa.PrivateKey
	// Serializer 序列化
	codec.Serializer
	// Deserializer 反序列化
	codec.Deserializer
	// Sign 签名
	Sign(hash []byte) (*Signature, error)
}

// Signature 签名交易
type Signature interface {
	// Sign 签名
	Sign(data []byte) (Signature, error)
	// Verify 验证签名
	Verify(pubKey *ecdsa.PublicKey, data []byte, signature []byte) (bool, error)
	// Bytes 获取签名数据
	Bytes() ([]byte, error)
}
