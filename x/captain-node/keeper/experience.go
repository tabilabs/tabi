package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captain-node/types"
)

func (k Keeper) UpdateAllUserExperience(
	ctx sdk.Context,
	userExperiences []*types.UserExperience,
) []sdk.Event {
	events := make([]sdk.Event, 0)
	for _, userExperience := range userExperiences {
		// update the experience of the user
		receiver := sdk.MustAccAddressFromBech32(userExperience.Receiver)
		oldExperience := k.GetExperience(ctx, receiver)
		k.incrExperience(ctx, receiver, userExperience.Experience)
		newExperience := k.GetExperience(ctx, receiver)
		events = append(
			events,
			sdk.NewEvent(
				types.EventTypeUpdateUserExperience,
				sdk.NewAttribute(types.AttributeKeyOwner, userExperience.Receiver),
				sdk.NewAttribute(types.AttributeKeyOldExperience, fmt.Sprintf("%d", oldExperience)),
				sdk.NewAttribute(types.AttributeKeyNewExperience, fmt.Sprintf("%d", newExperience)),
			),
		)
	}

	return nil
}

func (k Keeper) GetExperience(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	return k.getExperience(ctx, owner)
}

func (k Keeper) incrExperience(ctx sdk.Context, owner sdk.AccAddress, experience uint64) {
	store := ctx.KVStore(k.storeKey)

	// get the experience
	result := k.getExperience(ctx, owner) + experience
	store.Set(types.ExperienceStoreKey(owner), sdk.Uint64ToBigEndian(result))
}

func (k Keeper) decrExperience(ctx sdk.Context, owner sdk.AccAddress, experience uint64) {
	store := ctx.KVStore(k.storeKey)

	// get the experience
	result := k.getExperience(ctx, owner) - experience
	store.Set(types.ExperienceStoreKey(owner), sdk.Uint64ToBigEndian(result))
}

func (k Keeper) getExperience(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ExperienceStoreKey(owner))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
