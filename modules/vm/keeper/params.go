package keeper

import (
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	DefaultParamspace = types.ModuleName
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) GetMaxCodeSize(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxCodeSize, &res)
	return
}

func (k Keeper) SetMaxCodeSize(ctx sdk.Context, maxCodeSize uint64) {
	k.paramstore.Set(ctx, types.KeyMaxCodeSize, maxCodeSize)
}

func (k Keeper) GetVMFeeParams(ctx sdk.Context) (params [256]uint64) {
	k.paramstore.Get(ctx, types.KeyVMFeeParams, &params)
	return
}

func (k Keeper) SetVMFeeParams(ctx sdk.Context, params [256]uint64) {
	k.paramstore.Set(ctx, types.KeyVMFeeParams, params)
}

func (k Keeper) GetParams(ctx sdk.Context) (res types.Params) {
	return types.NewParams(
		k.GetMaxCodeSize(ctx),
		k.GetVMFeeParams(ctx),
	)
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
