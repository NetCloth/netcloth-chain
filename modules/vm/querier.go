package vm

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/vm/keeper"
	"github.com/netcloth/netcloth-chain/modules/vm/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func NewQuerier(k keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryCode:
			return queryCode(ctx, req, k)
		case types.QueryState:
			return queryState(ctx, req, k)
		case types.QueryStorage:
			return queryStorage(ctx, path, k)
		case types.QueryTxLogs:
			return queryTxLogs(ctx, path, k)
		case types.QueryCreateFee:
			return queryCreateFee(ctx, req, k)
		case types.QueryCallFee:
			return queryCallFee(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown vm query endpoint")
		}
	}
}

func queryParameters(ctx sdk.Context, k keeper.Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryCode(ctx sdk.Context, req abci.RequestQuery, k keeper.Keeper) ([]byte, sdk.Error) {
	if len(req.Data) != 20 {
		return nil, sdk.ErrInvalidAddress("address invalid")
	}

	accAddr := sdk.AccAddress(req.Data)
	code := k.GetCode(ctx, accAddr)

	return code, nil
}

func queryState(ctx sdk.Context, req abci.RequestQuery, k keeper.Keeper) ([]byte, sdk.Error) {
	var p types.QueryContractStateParams
	codec.Cdc.UnmarshalJSON(req.Data, &p)

	st := StateTransition{
		Sender:    p.From,
		Recipient: p.To,
		Price:     sdk.NewInt(1000000),
		GasLimit:  10000000,
		Amount:    sdk.NewInt(0),
		Payload:   p.Data,
		StateDB:   types.NewStateDB(k.StateDB).WithContext(ctx),
	}

	_, result := st.TransitionCSDB(ctx)

	return result.Data, nil
}

func queryStorage(ctx sdk.Context, path []string, keeper keeper.Keeper) ([]byte, sdk.Error) {
	addr, _ := sdk.AccAddressFromBech32(path[1])
	key := sdk.HexToHash(path[2])
	val := keeper.GetState(ctx, addr, key)
	bRes := types.QueryResStorage{Value: val.Bytes()}
	res, err := codec.MarshalJSONIndent(keeper.Cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}
	return res, nil
}

func queryTxLogs(ctx sdk.Context, path []string, keeper keeper.Keeper) ([]byte, sdk.Error) {
	txHash := sdk.HexToHash(path[1])
	logs := keeper.GetLogs(ctx, txHash)

	bRes := types.QueryLogs{Logs: logs}
	res, err := codec.MarshalJSONIndent(keeper.Cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryCreateFee(ctx sdk.Context, req abci.RequestQuery, k keeper.Keeper) ([]byte, sdk.Error) {
	var p types.QueryFeeParams
	codec.Cdc.UnmarshalJSON(req.Data, &p)

	st := StateTransition{
		Sender:    p.From,
		Recipient: nil,
		Price:     sdk.NewInt(100000000),
		GasLimit:  100000000,
		Amount:    sdk.NewInt(0),
		Payload:   p.Data,
		StateDB:   types.NewStateDB(k.StateDB).WithContext(ctx),
	}

	_, result := st.TransitionCSDB(ctx)

	if result.IsOK() {
		bRes := types.FeeResult{V: result.GasUsed}
		res, err := codec.MarshalJSONIndent(k.Cdc, bRes)
		if err != nil {
			panic("could not marshal result to JSON: " + err.Error())
		}
		return res, nil
	}

	return nil, sdk.ErrInternal("contract deploy failed")
}

func queryCallFee(ctx sdk.Context, req abci.RequestQuery, k keeper.Keeper) ([]byte, sdk.Error) {
	var p types.QueryFeeParams
	codec.Cdc.UnmarshalJSON(req.Data, &p)

	st := StateTransition{
		Sender:    p.From,
		Recipient: p.To,
		Price:     sdk.NewInt(100000000),
		GasLimit:  100000000,
		Amount:    sdk.NewInt(0),
		Payload:   p.Data,
		StateDB:   types.NewStateDB(k.StateDB).WithContext(ctx),
	}

	_, result := st.TransitionCSDB(ctx)

	if result.IsOK() {
		bRes := types.FeeResult{V: result.GasUsed}
		res, err := codec.MarshalJSONIndent(k.Cdc, bRes)
		if err != nil {
			panic("could not marshal result to JSON: " + err.Error())
		}
		return res, nil
	}

	return nil, sdk.ErrInternal("contract call failed")
}
