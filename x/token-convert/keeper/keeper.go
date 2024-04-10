package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	// TODO: add keepers
}

func NewKeeper(cdc codec.BinaryCodec, key *storetypes.StoreKey) Keeper {
	panic("TODO: implement me")
}
