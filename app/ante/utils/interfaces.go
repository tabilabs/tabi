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

package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BankKeeper defines the exposed interface for using functionality of the bank keeper
// in the context of the AnteHandler utils package.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// DistributionKeeper defines the exposed interface for using functionality of the distribution
// keeper in the context of the AnteHandler utils package.
type DistributionKeeper interface {
	WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error)
}

// StakingKeeper defines the exposed interface for using functionality of the staking keeper
// in the context of the AnteHandler utils package.
type StakingKeeper interface {
	BondDenom(ctx sdk.Context) string
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress, fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
}

// CaptainsKeeper defines the exposed interface for using functionality of the captains keeper
// in the context of the AnteHandler utils package.
type CaptainsKeeper interface {
	IsStandByPhase(ctx sdk.Context) bool
}
