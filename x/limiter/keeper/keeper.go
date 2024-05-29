package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/limiter/types"
)

// Keeper of the limiter store
type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey

	authority sdk.AccAddress
}

// NewKeeper creates a new limiter Keeper instance
func NewKeeper(cdc codec.Codec, key storetypes.StoreKey, authority sdk.AccAddress) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		authority: authority,
	}
}

// SetParams sets the parameters of the limiter module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)
	return nil
}

// GetParams returns the total set of limiter parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if len(bz) == 0 {
		return types.Params{}
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// AddAllowListMember adds a member to the allow list
func (k Keeper) AddAllowListMember(ctx sdk.Context, member string) error {
	params := k.GetParams(ctx)
	if _, found := isInAllowList(member, params.AllowList); found {
		return types.ErrMemberAlreadyExisted
	}
	params.AllowList = append(params.AllowList, member)
	return k.SetParams(ctx, params)
}

// RemoveAllowListMember removes a member from the allow list
func (k Keeper) RemoveAllowListMember(ctx sdk.Context, member string) error {
	params := k.GetParams(ctx)

	// check if the allow list is empty
	if len(params.AllowList) == 0 {
		return types.ErrEmptyAllowList
	}

	// check if the member is in the allow list
	index, found := isInAllowList(member, params.AllowList)
	if !found {
		return types.ErrMemberNotFound
	}

	newAllowList := make([]string, len(params.AllowList)-1)
	copy(newAllowList, params.AllowList[:index])
	copy(newAllowList[index:], params.AllowList[index+1:])
	params.AllowList = newAllowList

	return k.SetParams(ctx, params)
}

// SetEnabled sets the enabled status of the limiter module
func (k Keeper) SetEnabled(ctx sdk.Context, enabled bool) error {
	params := k.GetParams(ctx)
	params.Enabled = enabled
	return k.SetParams(ctx, params)
}

// IsEnabled returns the enabled status of the limiter module
func (k Keeper) IsEnabled(ctx sdk.Context) bool {
	params := k.GetParams(ctx)
	return params.Enabled
}

// IsAuthorized checks if the addr is in white list.
func (k Keeper) IsAuthorized(ctx sdk.Context, addr sdk.AccAddress) bool {
	params := k.GetParams(ctx)
	_, found := isInAllowList(addr.String(), params.AllowList)
	return found
}

// isInAllowList checks if the addr is in the list.
func isInAllowList(addr string, list []string) (int, bool) {
	for i, member := range list {
		if member == addr {
			return i, true
		}
	}
	return 0, false
}
