package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captain-node/types"
)

// InitGenesis stores the NFT genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	if err := data.Validate(); err != nil {
		panic(fmt.Errorf("failed to initialize mint genesis state: %s", err.Error()))
	}

	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(fmt.Errorf("failed to set mint genesis state: %s", err.Error()))
	}

	for _, division := range data.Divisions {
		if err := k.SaveDivision(ctx, *division); err != nil {
			panic(fmt.Errorf("failed to save division: %s", err.Error()))
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// todo: Entry
	return &types.GenesisState{
		Params:    k.GetParams(ctx),
		Divisions: k.GetDivisions(ctx),
	}
}
