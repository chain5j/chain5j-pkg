package json

import (
	"io"
	"math/big"
	"strings"
	"unsafe"

	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/json-iterator/go"
)

func init() {
	jsoniter.RegisterTypeEncoder("[]uint8", &byteCodec{})
	jsoniter.RegisterTypeDecoder("[]uint8", &byteCodec{})
	jsoniter.RegisterTypeEncoder("big.Int", &bigIntCodec{})
	jsoniter.RegisterTypeDecoder("big.Int", &bigIntCodec{})
}

func Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsoniter.MarshalIndent(v, prefix, indent)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return jsoniter.NewDecoder(r)
}

type byteCodec struct{}

func (codec *byteCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	b := *((*[]byte)(ptr))

	stream.WriteString(hexutil.Encode(b))
}

func (codec *byteCodec) IsEmpty(ptr unsafe.Pointer) bool {
	b := *((*[]byte)(ptr))

	return len(b) == 0
}

func (codec *byteCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	hex := iter.ReadString()

	b, err := hexutil.Decode(hex)
	if err != nil {
		return
	}
	*((*[]byte)(ptr)) = b
}

type bigIntCodec struct{}

func (codec *bigIntCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	i := *((*big.Int)(ptr))

	stream.WriteString(hexutil.EncodeBig(&i))
}

func (codec *bigIntCodec) IsEmpty(ptr unsafe.Pointer) bool {

	return false
}

func (codec *bigIntCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	result := new(big.Int)

	any := iter.ReadAny()
	valueType := any.ValueType()
	switch valueType {
	case jsoniter.StringValue:
		str := any.ToString()
		if !strings.HasPrefix(str, "0x") {
			result.UnmarshalText([]byte(str))
			*((*big.Int)(ptr)) = *result
			return
		}
		result, err := hexutil.DecodeBig(str)
		if err != nil {
			return
		}

		*((*big.Int)(ptr)) = *result
		break
	case jsoniter.NumberValue:
		num := any.ToInt64()
		result = big.NewInt(num)
		*((*big.Int)(ptr)) = *result
		break
	}
}
