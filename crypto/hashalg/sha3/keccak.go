// Package sha3
//
// @author: xwc1125
package sha3

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Keccak512 calculates and returns the Keccak512 hash of the input data.
func Keccak512(data ...[]byte) []byte {
	d := NewKeccak512()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
