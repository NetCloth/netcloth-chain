package vm

import (
	"math/big"

	"github.com/netcloth/netcloth-chain/modules/vm/types"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/netcloth/netcloth-chain/types"
)

var emptyCodeHash = crypto.Sha256(nil)

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(sdk.AccAddress, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(sdk.AccAddress, sdk.AccAddress, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) sdk.Hash
)

func run(evm *EVM, contract *Contract, input []byte, readOnly bool) ([]byte, error) {
	return nil, ErrNoCompatibleInterpreter
}

type codeAndHash struct {
	code []byte
	hash sdk.Hash
}

func (c *codeAndHash) Hash() sdk.Hash {
	if c.hash == (sdk.Hash{}) {
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
	StateDB *CommitStateDB

	// depth is the current call stack
	depth int

	chainConfig *ChainConfig

	// virtual machine configuration options used to initialise the vm
	vmConfig Config

	interpreters []Interpreter
	interpreter  Interpreter

	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32

	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

func NewEVM(ctx Context, statedb CommitStateDB, vmConfig Config) *EVM {
	evm := &EVM{
		Context:      ctx,
		StateDB:      &statedb,
		vmConfig:     vmConfig,
		interpreters: make([]Interpreter, 0, 1),
	}

	//evm.interpreters = append(evm.interpreters, NewEVMInterpreter(evm, vmConfig))
	evm.interpreter = evm.interpreters[0]
	return evm
}

// Interpreter returns the current interpreter
func (evm *EVM) Interpreter() Interpreter {
	return evm.interpreter
}

// Create creates a new contract using code as deployment code
func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr sdk.Address, leftOverGas uint64, err error) {
	//contractAddr = CreateAddress(caller.Address(), evm.StateDB.GetNonce(caller.Address()))
	//return evm.create(caller, &codeAndHash{code: code}, gas, value, contractAddr)
	return nil, nil, 0, nil
}

func (evm *EVM) create(caller ContractRef, codeAndHash *codeAndHash, gas uint64, value *big.Int, address sdk.AccAddress) ([]byte, sdk.AccAddress, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the limit
	if evm.depth > int(types.CallCreateDepth) {
		return nil, sdk.AccAddress{}, gas, ErrDepth
	}

	if !evm.CanTransfer(caller.Address(), value) {
		return nil, sdk.AccAddress{}, gas, ErrInsufficientBalance
	}

	nonce := evm.StateDB.GetNonce(caller.Address())
	evm.StateDB.SetNonce(caller.Address(), nonce+1)

	// Ensure there's no existing contract already at the designated address
	//contractHash := evm.StateDB.GetCodeHash(caller.Address())
	//if evm.StateDB.GetNonce(address) != 0 || (contractHash != (sdk.Hash{})) {
	//	return nil, sdk.AccAddress{}, 0, ErrContractAddressCollision
	//}

	// Create a new account on the state
	//snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(address)
	evm.StateDB.SetNonce(address, 1)
	evm.Transfer(caller.Address(), address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NetContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, address, gas, nil
	}

	//start := time.Now()
	ret, err := run(evm, contract, nil, false)

	maxCodeSizeExceeded := len(ret) > MaxCodeSize
	if err == nil && !maxCodeSizeExceeded {
		createGas := uint64(len(ret)) * types.CreateAccountGas
		if contract.UseGas(createGas) {
			evm.StateDB.SetCode(address, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (err != ErrCodeStoreOutOfGas)) {
		//evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	if evm.vmConfig.Debug && evm.depth == 0 {
		//evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}
	return ret, address, contract.Gas, err
}

// Create2 creates a new contract using code as deployment code.
//
// The different between Create2 with Create is Create2 uses sha3(0xff ++ msg.sender ++ salt ++ sha3(init_code))[12:]
// instead of the usual sender-and-nonce-hash as the address where the contract is initialized at.
func (evm *EVM) Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int, salt *big.Int) (ret []byte, contractAddr sdk.AccAddress, leftOverGas uint64, err error) {
	// TODO
	return
}

func (evm *EVM) Call(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	// TODO
	return
}

func (evm *EVM) CallCode(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	// TODO
	return
}

func (evm *EVM) DelegateCall(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	// TODO
	return
}

func (evm *EVM) StaticCall(caller ContractRef, addr sdk.AccAddress, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	// TODO
	return
}