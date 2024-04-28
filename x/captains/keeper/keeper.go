package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

type Keeper struct {
	cdc        codec.Codec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	stakingKeeper types.StakingKeeper

	authority sdk.AccAddress
}

func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	authority sdk.AccAddress,
) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramSpace,
		authority:  authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}
