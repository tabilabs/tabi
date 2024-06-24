package types // noalias

import (
	captainnodetypes "github.com/tabilabs/tabi/x/captains/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

	GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []captainnodetypes.Node)

	// GetCurrentEpoch return the current epoch id.
	GetCurrentEpoch(ctx sdk.Context) uint64

	// CalcNodeCumulativeEmissionByEpoch returns the cumulative emission of the node by the end of epoch.
	CalcNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec

	// GetNodeClaimedEmission returns the historical emission the last time user claimed.
	GetNodeClaimedEmission(ctx sdk.Context, nodeID string) sdk.Dec

	// UpdateGlobalAndNodeClaimedEmission updates the node_historical_emission_on_last_claim.
	UpdateGlobalAndNodeClaimedEmission(ctx sdk.Context, nodeID string) error
}

type MintKeeper interface {
	GetMinter(ctx sdk.Context) minttypes.Minter
}
