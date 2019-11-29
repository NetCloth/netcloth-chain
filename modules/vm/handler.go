package vm

import (
	"fmt"

	"github.com/netcloth/netcloth-chain/modules/vm/keeper"
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
	ctx.Logger().Info("handleMsgContractCreate ...")
	err := msg.ValidateBasic()
	if err != nil {
		return err.Result()
	}

	err = k.DoContractCreate(ctx, msg)
	if err != nil {
		return err.Result()
	}

	// check recur deep < 1024

	st := StateTransition{}
	_, res := st.TransitionCSDB(ctx)
	return res
}

func handleMsgContractCall(ctx sdk.Context, msg MsgContractCall, k Keeper) sdk.Result {
	ctx.Logger().Info("handleMsgContractCall ...")

	st := StateTransition{}
	_, res := st.TransitionCSDB(ctx)
	return res
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}