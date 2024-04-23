package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

func (k Keeper) SetAuthorizedMembers(ctx sdk.Context, members []string) ([]sdk.Event, error) {
	params := k.GetParams(ctx)

	events := make([]sdk.Event, 0)
	for _, member := range members {
		allowAdd := true
		for _, authzMember := range params.AuthorizedMembers {
			if authzMember == member {
				allowAdd = false
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
		return nil, err
	}

	return events, nil
}

func (k Keeper) RemoveCaller(ctx sdk.Context, members []string) ([]sdk.Event, error) {
	params := k.GetParams(ctx)

	events := make([]sdk.Event, 0)
	for _, member := range members {
		allowRemove := false
		for i, authzMember := range params.AuthorizedMembers {
			if authzMember == member {
				params.AuthorizedMembers = append(params.AuthorizedMembers[:i], params.AuthorizedMembers[i+1:]...)
				allowRemove = true
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
	if err := k.SetParams(ctx, params); err != nil {
		return nil, err
	}

	return events, nil
}

// HasAuthorizedMember asserts whether a Callers is already registered
func (k Keeper) HasAuthorizedMember(ctx sdk.Context, member sdk.AccAddress) bool {
	params := k.GetParams(ctx)

	allowAuthz := false
	for _, authzMember := range params.AuthorizedMembers {
		if sdk.MustAccAddressFromBech32(authzMember).Equals(member) {
			allowAuthz = true
		}
	}
	return allowAuthz
}

func (k Keeper) GetAuthorizedMembers(ctx sdk.Context) []string {
	return k.GetParams(ctx).AuthorizedMembers
}
