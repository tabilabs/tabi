// Copyright 2024 Tabi Foundation
// This file is part of the Tabi Network packages.
//
// Tabi is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Tabi packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

package claims

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabi/tabi/v1/x/claims/keeper"
	"github.com/tabi/tabi/v1/x/claims/types"
)

// InitGenesis initializes the claim module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	totalEscrowed := sdk.ZeroInt()
	sumUnclaimed := sdk.ZeroInt()
	numActions := sdk.NewInt(4)

	// ensure claim module account is set on genesis
	if acc := k.GetModuleAccount(ctx); acc == nil {
		panic("the claim module account has not been set")
	}

	// set the start time to the current block time by default
	if data.Params.AirdropStartTime.IsZero() {
		data.Params.AirdropStartTime = ctx.BlockTime()
	}

	err := k.SetParams(ctx, data.Params)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error setting params"))
	}

	escrowedCoins := k.GetModuleAccountBalances(ctx)
	if escrowedCoins != nil {
		totalEscrowed = escrowedCoins.AmountOfNoDenomValidation(data.Params.ClaimsDenom)
	}

	for _, claimsRecord := range data.ClaimsRecords {
		addr := sdk.MustAccAddressFromBech32(claimsRecord.Address)
		cr := types.ClaimsRecord{
			InitialClaimableAmount: claimsRecord.InitialClaimableAmount,
			ActionsCompleted:       claimsRecord.ActionsCompleted,
		}

		if len(cr.ActionsCompleted) != len(types.Action_name)-1 {
			panic(fmt.Errorf("invalid actions completed length for address %s", claimsRecord.Address))
		}

		initialClaimablePerAction := claimsRecord.InitialClaimableAmount.Quo(numActions)

		for _, actionCompleted := range cr.ActionsCompleted {
			if !actionCompleted {
				// NOTE: only add the initial claimable amount per action for the ones that haven't been claimed
				sumUnclaimed = sumUnclaimed.Add(initialClaimablePerAction)
			}
		}

		k.SetClaimsRecord(ctx, addr, cr)
	}

	// check for equal only for unclaimed actions
	if !sumUnclaimed.Equal(totalEscrowed) {
		panic(
			fmt.Errorf(
				"sum of unclaimed amount ≠ escrowed module account amount (%s ≠ %s)",
				sumUnclaimed, totalEscrowed,
			),
		)
	}
}

// ExportGenesis returns the claim module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:        k.GetParams(ctx),
		ClaimsRecords: k.GetClaimsRecords(ctx),
	}
}
