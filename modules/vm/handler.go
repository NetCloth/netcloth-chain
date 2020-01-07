package vm

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/modules/vm/keeper"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgContractCreate:
			return handleMsgContractCreate(ctx, msg, k)
		case MsgContractCall:
			return handleMsgContractCall(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("Unrecognized Msg type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgContractCreate(ctx sdk.Context, msg MsgContractCreate, k Keeper) sdk.Result {
	// validate msg
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	st := StateTransition{
		Sender:    msg.From,
		Recipient: nil,
		Price:     sdk.NewInt(1000000),
		GasLimit:  10000000,
		Amount:    msg.Amount.Amount,
		Payload:   msg.Code,
		StateDB:   k.StateDB.WithContext(ctx),
	}

	f := k.GetVMOpGasParams(ctx)
	f1 := k.GetVMCommonGasParams(ctx)
	_, res := st.TransitionCSDB(ctx, &f, &f1)
	if !res.IsOK() {
		// return vm error
		return res
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Data: res.Data, GasUsed: res.GasUsed, Events: ctx.EventManager().Events()}
}

func handleMsgContractCall(ctx sdk.Context, msg MsgContractCall, k Keeper) sdk.Result {
	// validate msg
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	// check code
	if code := k.GetCode(ctx, msg.Recipient); code == nil {
		return ErrNoCodeExist().Result()
	}

	st := StateTransition{
		Sender:    msg.From,
		Recipient: msg.Recipient,
		Price:     sdk.NewInt(1000000),
		GasLimit:  10000000,
		Payload:   msg.Payload,
		Amount:    msg.Amount.Amount,
		StateDB:   k.StateDB.WithContext(ctx),
	}

	f := k.GetVMOpGasParams(ctx)
	f1 := k.GetVMCommonGasParams(ctx)
	_, res := st.TransitionCSDB(ctx, &f, &f1)
	if !res.IsOK() {
		// return vm error
		return res
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Data: res.Data, GasUsed: res.GasUsed, Events: ctx.EventManager().Events()}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	k.StateDB.WithContext(ctx).Commit(true)
	return []abci.ValidatorUpdate{}
}
