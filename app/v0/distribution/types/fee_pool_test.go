package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/netcloth/netcloth-chain/types"
)

func TestValidateGenesis(t *testing.T) {

	fp := InitialFeePool()
	require.Nil(t, fp.ValidateGenesis())

	fp2 := FeePool{CommunityPool: sdk.DecCoins{{Denom: "nch", Amount: sdk.NewDec(-1)}}}
	require.NotNil(t, fp2.ValidateGenesis())

}
