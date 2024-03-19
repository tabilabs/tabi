package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captain-node/types"
)

func (k Keeper) UpdateAllNodesPowerOnPeriod(
	ctx sdk.Context,
	NodesPowerOnPeriod []*types.CaptainNodePowerOnPeriod,
) []sdk.Event {
	events := make([]sdk.Event, 0)
	params := k.GetParams(ctx)
	for _, nodePowerOnPeriod := range NodesPowerOnPeriod {
		// update the power on period of the node
		oldPowerOnPeriod := k.GetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId)
		powerOnPeriod := calculatePowerOnPeriod(nodePowerOnPeriod.PowerOnPeriod, params.MaximumPowerOnPeriod)
		k.SetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId, powerOnPeriod)
		newPowerOnPeriod := k.GetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId)
		events = append(
			events,
			sdk.NewEvent(
				types.EventTypeUpdatePowerOnPeriod,
				sdk.NewAttribute(types.AttributeKeyNodeID, nodePowerOnPeriod.NodeId),
				sdk.NewAttribute(types.AttributeKeyOldPowerOnPeriod, oldPowerOnPeriod.String()),
				sdk.NewAttribute(types.AttributeKeyNewPowerOnPeriod, newPowerOnPeriod.String()),
			),
		)
	}
	return events
}

// SetNodePowerOnPeriod set the proportion of online nodes
func (k Keeper) SetNodePowerOnPeriod(ctx sdk.Context, nodeID string, powerOnPeriod sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.NodePowerOnPeriodStoreKey(nodeID), powerOnPeriod.BigInt().Bytes())
}

// GetNodePowerOnPeriod returns the proportion of online nodes
func (k Keeper) GetNodePowerOnPeriod(ctx sdk.Context, nodeID string) sdk.Dec {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.NodePowerOnPeriodStoreKey(nodeID))
	if bz == nil {
		return sdk.ZeroDec()
	}

	return sdk.MustNewDecFromStr(string(bz))
}
