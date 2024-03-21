package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// expected keepers for captain-node module

type StakeKeeper interface {
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation stakingtypes.DelegationI) (stop bool))
	BondDenom(ctx sdk.Context) (res string)
}
