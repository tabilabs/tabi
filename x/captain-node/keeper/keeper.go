package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc        codec.Codec
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
	authority sdk.AccAddress
}
