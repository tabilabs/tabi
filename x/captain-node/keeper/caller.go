package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captain-node/types"
)

func (k Keeper) SetCaller(ctx sdk.Context, callerAddrs []string) ([]sdk.Event, error) {
	params := k.GetParams(ctx)

	events := make([]sdk.Event, 0)
	for _, callerAddr := range callerAddrs {
		allowAdd := true
		for _, caller := range params.Callers {
			if caller == callerAddr {
				allowAdd = false
			}
		}
		if allowAdd {
			params.Callers = append(params.Callers, callerAddr)
			events = append(
				events,
				sdk.NewEvent(
					types.EventTypeAddCaller,
					sdk.NewAttribute(types.AttributeCaller, callerAddr),
				),
			)
		}
	}
	if err := k.SetParams(ctx, params); err != nil {
		return nil, err
	}

	return events, nil
}

func (k Keeper) RemoveCaller(ctx sdk.Context, callerAddrs []string) ([]sdk.Event, error) {
	params := k.GetParams(ctx)

	events := make([]sdk.Event, 0)
	for _, callerAddr := range callerAddrs {
		allowRemove := false
		for i, caller := range params.Callers {
			if caller == callerAddr {
				params.Callers = append(params.Callers[:i], params.Callers[i+1:]...)
				allowRemove = true
			}
		}
		if allowRemove {
			events = append(
				events,
				sdk.NewEvent(
					types.EventTypeRemoveCaller,
					sdk.NewAttribute(types.AttributeCaller, callerAddr),
				),
			)
		}
	}
	if err := k.SetParams(ctx, params); err != nil {
		return nil, err
	}

	return events, nil
}

// AuthCaller asserts whether a Callers is already registered
func (k Keeper) AuthCaller(ctx sdk.Context, callerAddr sdk.AccAddress) bool {
	params := k.GetParams(ctx)

	allowCall := false
	for _, caller := range params.Callers {
		if sdk.MustAccAddressFromBech32(caller).Equals(callerAddr) {
			allowCall = true
		}
	}
	return allowCall
}

func (k Keeper) GetCallers(ctx sdk.Context) []string {
	return k.GetParams(ctx).Callers
}
