// description: chain5j-core 
// 
// @author: xwc1125
// @date: 2019/9/12
package types

//// 地址的长度
//const HexAddressLength = 20
//
//var addressT = reflect.TypeOf(HexAddress{})
//
//// HexAddress represents the 20 byte address of an  account.
//type HexAddress [HexAddressLength]byte
//
//var (
//	_               Address = HexAddress{}
//	EmptyHexAddress         = HexAddress{}
//)
//
//func NewHexAddress() HexAddress {
//	return HexAddress{}
//}
//
//func (a HexAddress) Len() int {
//	return HexAddressLength
//}
//
//// Bytes gets the string representation of the underlying address.
//func (a HexAddress) Bytes() []byte { return a[:] }
//
//// String implements fmt.Stringer.
//func (a HexAddress) String() string {
//	return a.Hex()
//}
//
//func (a HexAddress) FromBytes(bytes []byte) (Address, error) {
//	var a1 HexAddress
//	a1.SetBytes(bytes)
//	return a1, nil
//}
//
//func (a HexAddress) FromStr(addr string) (Address, error) {
//	return a.FromBytes(hexutil.FromHex(addr))
//}
//
//func (a HexAddress) Validate(addr string) bool {
//	if hexutil.HasHexPrefix(addr) {
//		addr = addr[2:]
//	}
//	return len(addr) == 2*HexAddressLength && hexutil.IsHex(addr)
//}
//
//func (a HexAddress) Nil() bool {
//	return a == HexAddress{}
//}
//
//// Hash converts an address to a hash by left-padding it with zeros.
//func (a HexAddress) Hash() Hash { return BytesToHash(a[:]) }
//
//// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
//// without going through the stringer interface used for logging.
//func (a HexAddress) Format(s fmt.State, c rune) {
//	fmt.Fprintf(s, "%"+string(c), a[:])
//}
//
//// MarshalText returns the hex representation of a.
//func (a HexAddress) MarshalText() ([]byte, error) {
//	return hexutil.Bytes(a[:]).MarshalText()
//}
//
//// UnmarshalText parses a hash in hex syntax.
//func (a *HexAddress) UnmarshalText(input []byte) error {
//	return hexutil.UnmarshalFixedText("HexAddress", input, a[:])
//}
//
//// UnmarshalJSON parses a hash in hex syntax.
//func (a *HexAddress) UnmarshalJSON(input []byte) error {
//	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
//}
//
//// BigToAddress returns Address with byte values of b.
//// If b is larger than len(h), b will be cropped from the left.
//func BigToAddress(b *big.Int) (Address, error) {
//	var a Address
//	return a.FromBytes(b.Bytes())
//}
//
//// HexToAddress returns Address with byte values of s.
//// If s is larger than len(h), s will be cropped from the left.
//func HexToAddress(s string) (Address, error) {
//	var a Address
//	return a.FromBytes(hexutil.FromHex(s))
//}
//
//// Big converts an address to a big integer.
//func (a HexAddress) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }
//
//// Hex returns an EIP55-compliant hex string representation of the address.
//func (a HexAddress) Hex() string {
//	unchecksummed := hex.EncodeToString(a[:])
//	sha := sha3.NewKeccak256()
//	sha.Write([]byte(unchecksummed))
//	hash := sha.Sum(nil)
//
//	result := []byte(unchecksummed)
//	for i := 0; i < len(result); i++ {
//		hashByte := hash[i/2]
//		if i%2 == 0 {
//			hashByte = hashByte >> 4
//		} else {
//			hashByte &= 0xf
//		}
//		if result[i] > '9' && hashByte > 7 {
//			result[i] -= 32
//		}
//	}
//	return "0x" + string(result)
//}
//
//// SetBytes sets the address to the value of b.
//// If b is larger than len(a) it will panic.
//func (a *HexAddress) SetBytes(b []byte) {
//	if len(b) > len(a) {
//		b = b[len(b)-HexAddressLength:]
//	}
//	copy(a[HexAddressLength-len(b):], b)
//}
//
//// Scan implements Scanner for database/sql.
//func (a *HexAddress) Scan(src interface{}) error {
//	srcB, ok := src.([]byte)
//	if !ok {
//		return fmt.Errorf("can't scan %T into HexAddress", src)
//	}
//	if len(srcB) != HexAddressLength {
//		return fmt.Errorf("can't scan []byte of len %d into HexAddress, want %d", len(srcB), HexAddressLength)
//	}
//	copy(a[:], srcB)
//	return nil
//}
//
//// Value implements valuer for database/sql.
//func (a HexAddress) Value() (driver.Value, error) {
//	return a[:], nil
//}
