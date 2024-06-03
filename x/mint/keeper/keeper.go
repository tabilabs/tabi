package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/tabilabs/tabi/x/mint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// keeper of the mint store
type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	distrKeeper   types.DistrKeeper

	feeCollectorName string

	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority sdk.AccAddress
}

// NewKeeper returns a mint keeper
func NewKeeper(cdc codec.Codec, authority sdk.AccAddress, key storetypes.StoreKey,
	ak types.AccountKeeper, bk types.BankKeeper, dk types.DistrKeeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	keeper := Keeper{
		storeKey:         key,
		cdc:              cdc,
		accountKeeper:    ak,
		bankKeeper:       bk,
		distrKeeper:      dk,
		feeCollectorName: feeCollectorName,
		authority:        authority,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("%s", types.ModuleName))
}
