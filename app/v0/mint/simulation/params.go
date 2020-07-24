package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/netcloth/netcloth-chain/app/v0/mint/internal/types"
	"github.com/netcloth/netcloth-chain/app/v0/simulation"
	simtypes "github.com/netcloth/netcloth-chain/types/simulation"
)

const (
	keyInflationRateChange = "InflationRateChange"
	keyInflationMax        = "InflationMax"
	keyInflationMin        = "InflationMin"
	keyGoalBonded          = "GoalBonded"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyInflationRateChange,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenInflationRateChange(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyInflationMax,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenInflationMax(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyInflationMin,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenInflationMin(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyGoalBonded,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenGoalBonded(r))
			},
		),
	}
}
