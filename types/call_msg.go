package types

import "math/big"

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	From     Address  // the sender of the 'transaction'
	To       *Address // the destination contract (nil for contract creation)
	Gas      uint64   // if 0, the call executes with near-infinite gas
	GasPrice *big.Int // wei <-> gas exchange ratio
	Value    *big.Int // amount of wei sent along with the call
	Data     []byte   // input data, usually an ABI-encoded contract method invocation
}
