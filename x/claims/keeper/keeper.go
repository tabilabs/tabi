package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tabilabs/tabi/x/claims/types"
)

type Keeper struct {
	cdc        codec.Codec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	// cosmos keepers
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper

	// self module keepers
	captainsKeeper types.CaptainsKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority sdk.AccAddress
}

// NewKeeper returns a mint keeper
func NewKeeper(cdc codec.Codec, authority sdk.AccAddress,
	key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper, ck types.CaptainsKeeper,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the climas collector account has not been set")
	}

	keeper := Keeper{
		storeKey:       key,
		cdc:            cdc,
		paramSpace:     paramSpace.WithKeyTable(types.ParamKeyTable()),
		authKeeper:     ak,
		bankKeeper:     bk,
		captainsKeeper: ck,
		authority:      authority,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}

// SetCaptainsKeeper
func (k *Keeper) SetCaptainsKeeper(ck types.CaptainsKeeper) {
	k.captainsKeeper = ck
}
