package json

import (
	"math/big"
	"testing"
)

func TestByteCodec(t *testing.T) {
	type tmpst struct {
		Data []byte
		Num  *big.Int
	}

	tmp := &tmpst{
		Data: []byte{0x1, 0x2, 0x3, 0xf},
		Num:  big.NewInt(100),
	}

	b, err := Marshal(tmp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("encode byte: %s", string(b))

	var dec tmpst
	if err := Unmarshal(b, &dec); err != nil {
		t.Fatal(err)
	}

	return
}

func TestBigIntCodec(t *testing.T) {
	type tmpst struct {
		Num *big.Int
	}

	tmp := &tmpst{
		Num: big.NewInt(100),
	}

	b, err := Marshal(tmp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("encode: %s", string(b))

	var dec tmpst
	if err := Unmarshal(b, &dec); err != nil {
		t.Fatal(err)
	}

	tmp2 := &tmpst{
	}

	b2, err := Marshal(tmp2)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("encode: %s", string(b2))

	if err := Unmarshal(b2, &dec); err != nil {
		t.Fatal(err)
	}

	return
}
