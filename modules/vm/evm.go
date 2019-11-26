package vm

import (
	"math/big"

	"github.com/tendermint/tendermint/crypto"

	"github.com/netcloth/netcloth-chain/modules/vm/common"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(sdk.AccAddress, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(sdk.AccAddress, sdk.AccAddress, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

type codeAndHash struct {
	code []byte
	hash common.Hash
}

func (c *codeAndHash) Hash() common.Hash {
	if c.hash == (common.Hash{}) {
		copy(c.hash[:], crypto.Sha256(c.code))
	}
	return c.hash
}

// Context provides the VM with auxiliary information.
// Once provided it shouldn't be modified
type Context struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Msg information
	Origin   sdk.AccAddress
	GasPrice *big.Int

	// Block information
	CoinBase    sdk.AccAddress
	GasLimit    uint64
	BlockNumber *big.Int
	Time        *big.Int
}

type EVM struct {
	Context

	// StateDB gives access to the underlying state
	StateDB StateDB

	// depth is the current call stack
	depth int

	chainConfig *ChainConfig

	// virtual machine configuration options used to initialise the vm
	vmConfig Config

	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32

	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

// Create creates a new contract using code as deployment code
func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr sdk.Address, leftOverGas uint64, err error) {
	// TODO
}
