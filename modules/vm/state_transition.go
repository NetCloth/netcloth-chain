package vm

import (
	"fmt"
	"math/big"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

// StateTransition defines data to transitionDB in vm
type StateTransition struct {
	Sender    sdk.AccAddress
	GasLimit  uint64
	Recipient sdk.AccAddress
	Amount    sdk.Int
	Payload   []byte
	StateDB   *types.CommitStateDB
}

func (st StateTransition) CanTransfer(acc sdk.AccAddress, amount *big.Int) bool {
	return st.StateDB.GetBalance(acc).Cmp(amount) >= 0
}

func (st StateTransition) Transfer(from, to sdk.AccAddress, amount *big.Int) {
	st.StateDB.SubBalance(from, amount)
	st.StateDB.AddBalance(to, amount)
}

func (st StateTransition) GetHashFn(header abci.Header) func(n uint64) sdk.Hash {
	return func(n uint64) sdk.Hash {
		var res = sdk.Hash{}
		blockID := header.GetLastBlockId()
		res.SetBytes(blockID.GetHash())
		return res
	}
}

func (st StateTransition) TransitionCSDB(ctx sdk.Context, constGasConfig *[256]uint64, vmCommonGasConfig *types.VMCommonGasParams) (*big.Int, sdk.Result) {
	st.StateDB.UpdateAccounts()
	evmCtx := Context{
		CanTransfer: st.CanTransfer,
		Transfer:    st.Transfer,
		GetHash:     st.GetHashFn(ctx.BlockHeader()),

		Origin: st.Sender,

		CoinBase:    ctx.BlockHeader().ProposerAddress, // TODO: should be proposer account address
		GasLimit:    st.GasLimit,
		BlockNumber: sdk.NewInt(ctx.BlockHeader().Height).BigInt(),
	}

	cfg := Config{OpConstGasConfig: constGasConfig, CommonGasConfig: vmCommonGasConfig}

	currentGasMeter := ctx.GasMeter()
	evm := NewEVM(evmCtx, st.StateDB.WithTxHash(tmhash.Sum(ctx.TxBytes())).WithContext(ctx.WithGasMeter(sdk.NewInfiniteGasMeter())), cfg)

	var (
		ret         []byte
		leftOverGas uint64
		addr        sdk.AccAddress
		vmerr       error
	)

	if st.Recipient.Empty() {
		ret, addr, leftOverGas, vmerr = evm.Create(st.Sender, st.Payload, st.GasLimit, st.Amount.BigInt())
		fmt.Fprint(os.Stderr, "\n\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
		fmt.Fprint(os.Stderr, "+                                                                             +\n")
		fmt.Fprint(os.Stderr, fmt.Sprintf("+         contractAddr = %s           +\n", addr))
		fmt.Fprint(os.Stderr, "+                                                                             +\n")
		fmt.Fprint(os.Stderr, "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n\n")
	} else {
		ret, leftOverGas, vmerr = evm.Call(st.Sender, st.Recipient, st.Payload, st.GasLimit, st.Amount.BigInt())
	}

	fmt.Fprint(os.Stderr, fmt.Sprintf("ret = %x, \nconsumed gas = %v , leftOverGas = %v, err = %v\n", ret, st.GasLimit-leftOverGas, leftOverGas, vmerr))

	ctx.WithGasMeter(currentGasMeter).GasMeter().ConsumeGas(st.GasLimit-leftOverGas, "EVM execution consumption")
	if vmerr != nil {
		return nil, vmerr
	}

	st.StateDB.Finalise(true)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeNewContract,
			sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
		),
	})

	return nil, sdk.Result{Data: ret, GasUsed: st.GasLimit - leftOverGas}
}

func DoStateTransition(ctx sdk.Context, msg types.MsgContract, k Keeper, gasLimit uint64, readonly bool) (*big.Int, sdk.Result) {
	st := StateTransition{
		Sender:    msg.From,
		Recipient: msg.To,
		GasLimit:  gasLimit,
		Payload:   msg.Payload,
		Amount:    msg.Amount.Amount,
		StateDB:   k.StateDB.WithContext(ctx),
	}

	if readonly {
		st.StateDB = types.NewStateDB(k.StateDB).WithContext(ctx)
	}

	opGasConfig := k.GetVMOpGasParams(ctx)
	commonGasConfig := k.GetVMCommonGasParams(ctx)

	return st.TransitionCSDB(ctx, &opGasConfig, &commonGasConfig)
}
