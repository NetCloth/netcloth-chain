package cli

import (
	"encoding/json"
	"fmt"
	v0 "github.com/netcloth/netcloth-chain/app/v0"
	"os"

	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/server"
)

// Validate genesis command takes
func ValidateGenesisCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validate-genesis [file]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "validates the genesis file at the default location or at the location passed as an arg",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			// Load default if passed no args, otherwise load passed file
			var genesis string
			if len(args) == 0 {
				genesis = ctx.Config.GenesisFile()
			} else {
				genesis = args[0]
			}

			fmt.Fprintf(os.Stderr, "validating genesis file at %s\n", genesis)

			var genDoc *tmtypes.GenesisDoc
			if genDoc, err = tmtypes.GenesisDocFromFile(genesis); err != nil {
				return fmt.Errorf("error loading genesis doc from %s: %s", genesis, err.Error())
			}

			var genState map[string]json.RawMessage
			if err = cdc.UnmarshalJSON(genDoc.AppState, &genState); err != nil {
				return fmt.Errorf("error unmarshaling genesis doc %s: %s", genesis, err.Error())
			}

			if err = v0.ModuleBasics.ValidateGenesis(genState); err != nil {
				return fmt.Errorf("error validating genesis file %s: %s", genesis, err.Error())
			}

			// TODO test to make sure initchain doesn't panic

			fmt.Printf("File at %s is a valid genesis file\n", genesis)
			return nil
		},
	}
}