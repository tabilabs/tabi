package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captain-node/types"
)

func (k Keeper) UpdateExtractableComputingPowerForUsers(
	ctx sdk.Context,
	extractableComputingPowers []*types.ExtractableComputingPower,
) []sdk.Event {
	events := make([]sdk.Event, 0)
	for _, extractableComputingPower := range extractableComputingPowers {
		// update the experience of the user
		receiver := sdk.MustAccAddressFromBech32(extractableComputingPower.Owner)
		oldExperience := k.getExtractableComputingPower(ctx, receiver)
		//incrExtractableComputingPower is a function that increments the extractable computing power of the user
		k.incrExtractableComputingPower(ctx, receiver, extractableComputingPower.Amount)
		newExperience := k.getExtractableComputingPower(ctx, receiver)
		events = append(
			events,
			sdk.NewEvent(
				types.EventTypeUpdateUserExperience,
				sdk.NewAttribute(types.AttributeKeyOwner, extractableComputingPower.Owner),
				sdk.NewAttribute(types.AttributeKeyOldExperience, fmt.Sprintf("%d", oldExperience)),
				sdk.NewAttribute(types.AttributeKeyNewExperience, fmt.Sprintf("%d", newExperience)),
			),
		)
	}

	return nil
}

func (k Keeper) incrExtractableComputingPower(
	ctx sdk.Context,
	owner sdk.AccAddress,
	amount uint64,
) {
	store := ctx.KVStore(k.storeKey)

	// get the experience
	result := k.getExtractableComputingPower(ctx, owner) + amount
	store.Set(types.ExperienceStoreKey(owner), sdk.Uint64ToBigEndian(result))
}

func (k Keeper) decrExtractableComputingPower(
	ctx sdk.Context,
	owner sdk.AccAddress,
	amount uint64,
) {
	store := ctx.KVStore(k.storeKey)

	// get the experience
	result := k.getExtractableComputingPower(ctx, owner) - amount
	store.Set(types.ExperienceStoreKey(owner), sdk.Uint64ToBigEndian(result))
}

func (k Keeper) getExtractableComputingPower(ctx sdk.Context, owner sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ExperienceStoreKey(owner))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
