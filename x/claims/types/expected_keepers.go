package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	captainnodetypes "github.com/tabilabs/tabi/x/captain-node/types"
	minttypes "github.com/tabilabs/tabi/x/mint/types"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, types.ModuleAccountI)
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

// BankKeeper defines the contract needed to be fulfilled for banking and supply
// dependencies.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(
		ctx sdk.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToModule(
		ctx sdk.Context,
		senderModule, recipientModule string,
		amt sdk.Coins,
	) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error

	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
}

type StakingKeeper interface {
	GetDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress, maxRetrieve uint16) (delegations []stakingtypes.Delegation)
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	BondDenom(ctx sdk.Context) (res string)
	GetParams(ctx sdk.Context) stakingtypes.Params
}

// DistrKeeper defines the contract needed to be fulfilled for distribution keeper
type DistrKeeper interface {
	FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error
}

type CaptainsKeeper interface {
	GetParams(ctx sdk.Context) captainnodetypes.Params
	GetEpochs(ctx sdk.Context, sender sdk.Address) []string
	PruneEpochs(ctx sdk.Context, sender sdk.Address) // setHook

	GetUserHoldingQuantity(ctx sdk.Context, sender sdk.Address) uint64
}

type MintKeeper interface {
	GetDailyIssuance(ctx sdk.Context) sdk.Dec
	GetMinter(ctx sdk.Context) minttypes.Minter
}
