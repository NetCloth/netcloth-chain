package types

import (
	authexported "github.com/NetCloth/netcloth-chain/modules/auth/exported"
	sdk "github.com/NetCloth/netcloth-chain/types"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}