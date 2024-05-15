package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingKeeper interface {
	GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	BondDenom(ctx sdk.Context) (res string)
	GetParams(ctx sdk.Context) stakingtypes.Params
}
