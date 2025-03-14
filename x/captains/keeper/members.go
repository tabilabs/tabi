package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetAuthorizedMembers sets the list of authorized members
func (k Keeper) SetAuthorizedMembers(ctx sdk.Context, members []string) error {
	params := k.GetParams(ctx)
	events := make([]sdk.Event, 0)

	for _, member := range members {
		allowAdd := true
		for _, authzMember := range params.AuthorizedMembers {
			if authzMember == member {
				allowAdd = false
				break
			}
		}
		if allowAdd {
			params.AuthorizedMembers = append(params.AuthorizedMembers, member)
			events = append(
				events,
				sdk.NewEvent(
					types.EventTypeAddAuthorizedMembers,
					sdk.NewAttribute(types.AttributeKeyAuthorizedMember, member),
				),
			)
		}
	}

	if err := k.SetParams(ctx, params); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(events)

	return nil
}

// DeleteAuthorizedMembers deletes the list of authorized members
func (k Keeper) DeleteAuthorizedMembers(ctx sdk.Context, members []string) error {
	params := k.GetParams(ctx)

	events := make([]sdk.Event, 0)
	for _, member := range members {
		allowRemove := false
		for i, authzMember := range params.AuthorizedMembers {
			if authzMember == member {
				params.AuthorizedMembers = append(params.AuthorizedMembers[:i], params.AuthorizedMembers[i+1:]...)
				allowRemove = true
				break
			}
		}
		if allowRemove {
			events = append(
				events,
				sdk.NewEvent(
					types.EventTypeRemoveAuthorizedMembers,
					sdk.NewAttribute(types.AttributeKeyAuthorizedMember, member),
				),
			)
		}
	}

	if len(params.AuthorizedMembers) == 0 {
		return errorsmod.Wrap(types.ErrDeleteLastMember, "can not delete the last member")
	}

	if err := k.SetParams(ctx, params); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(events)

	return nil
}

// HasAuthorizedMember asserts whether an address is in authorized member lists.
func (k Keeper) HasAuthorizedMember(ctx sdk.Context, member sdk.AccAddress) bool {
	params := k.GetParams(ctx)

	allowAuthz := false
	for _, authzMember := range params.AuthorizedMembers {
		if sdk.MustAccAddressFromBech32(authzMember).Equals(member) {
			allowAuthz = true
			break
		}
	}
	return allowAuthz
}

// GetAuthorizedMembers returns the list of authorized members
func (k Keeper) GetAuthorizedMembers(ctx sdk.Context) []string {
	return k.GetParams(ctx).AuthorizedMembers
}
