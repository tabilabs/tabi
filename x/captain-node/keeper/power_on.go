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
	for _, nodePowerOnPeriod := range NodesPowerOnPeriod {
		// update the power on period of the node
		oldPowerOnPeriod := k.GetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId)
		k.SetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId, nodePowerOnPeriod.PowerOnPeriod)
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
func (k Keeper) SetNodePowerOnPeriod(ctx sdk.Context, nodeID string, proportion uint64) {
	store := ctx.KVStore(k.storeKey)

	powerOnPeriod := k.CalculatePowerOnPeriod(ctx, proportion)
	store.Set(types.NodePowerOnPeriodStoreKey(nodeID), powerOnPeriod.Bytes())
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
