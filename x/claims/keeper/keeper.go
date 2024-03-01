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

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority sdk.AccAddress
}

// NewKeeper returns a mint keeper
func NewKeeper(cdc codec.Codec, authority sdk.AccAddress,
	key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	ak types.AccountKeeper, bk types.BankKeeper) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ClaimsCollectorName); addr == nil {
		panic("the climas collector account has not been set")
	}

	keeper := Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace.WithKeyTable(types.ParamKeyTable()),
		accountKeeper: ak,
		bankKeeper:    bk,
		authority:     authority,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}
