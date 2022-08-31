package scrypt

import (
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	veryLightScryptN = 2
	veryLightScryptP = 1
)

func TestEncryptKey(t *testing.T) {
	prvKey, _ := hex.DecodeString("587ca4a15bc4d239cfba433dda03366506e99ecd2c529216eb3168b3e7806257")
	k := &Key{
		// Address:    addr,
		PrivateKey: prvKey,
	}
	keyJson, err := EncryptKey(k, "123456", veryLightScryptN, veryLightScryptP)
	if err != nil {
		t.Errorf("test: failed to recrypt key %v", err)
	}
	fmt.Println(string(keyJson))

	key, err := DecryptKey(keyJson, "123456")
	if err != nil {
		t.Errorf("test: failed to recrypt key %v", err)
	}
	prvKey1 := hex.EncodeToString(key.PrivateKey)
	fmt.Println("prvKey1", prvKey1)
}
