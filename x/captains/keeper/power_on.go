package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

var (
	smallestOperationalRate = sdk.NewDecWithPrec(47, 2) // constant 0.47
)

//func (k Keeper) UpdateAllNodesPowerOnPeriod(
//	ctx sdk.Context,
//	NodesPowerOnPeriod []*types.CaptainNodePowerOnPeriod,
//) []sdk.Event {
//	events := make([]sdk.Event, 0)
//
//	totalNumber := 0
//	for _, nodePowerOnPeriod := range NodesPowerOnPeriod {
//		// update the power on period of the node
//		oldPowerOnPeriod := k.GetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId)
//		if k.isAConsistentPowerOnPeriod(ctx, nodePowerOnPeriod.PowerOnPeriodRate) {
//			totalNumber += 1
//		}
//		// set the new power on period
//		k.SetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId, nodePowerOnPeriod.PowerOnPeriodRate)
//		newPowerOnPeriod := k.GetNodePowerOnPeriod(ctx, nodePowerOnPeriod.NodeId)
//		events = append(
//			events,
//			sdk.NewEvent(
//				types.EventTypeUpdatePowerOnPeriod,
//				sdk.NewAttribute(types.AttributeKeyNodeID, nodePowerOnPeriod.NodeId),
//				sdk.NewAttribute(types.AttributeKeyOldPowerOnPeriod, oldPowerOnPeriod.String()),
//				sdk.NewAttribute(types.AttributeKeyNewPowerOnPeriod, newPowerOnPeriod.String()),
//			),
//		)
//	}
//
//	operationalRate := k.calculateOperationalRate(ctx, totalNumber)
//	k.SetOperationalRate(ctx, operationalRate)
//
//	return events
//}

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

func (k Keeper) SetOperationalRate(ctx sdk.Context, operationalRate sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.OperationalRateKey, operationalRate.BigInt().Bytes())
}

func (k Keeper) GetOperationalRate(ctx sdk.Context) sdk.Dec {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.OperationalRateKey)
	if bz == nil {
		return smallestOperationalRate
	}

	return sdk.MustNewDecFromStr(string(bz))
}

func (k Keeper) calculateOperationalRate(ctx sdk.Context, totalNumber int) sdk.Dec {
	// Calculate OperationalRate = totalNumber / soldTotalCount
	soldTotalCount := k.GetNodeSequence(ctx) - 1
	operationalRate := sdk.NewDec(int64(totalNumber)).Quo(sdk.NewDec(int64(soldTotalCount)))
	// if operationalRate < 0.47 then operationalRate = 0.47
	// else operationalRate = operationalRate
	if operationalRate.LT(smallestOperationalRate) {
		operationalRate = smallestOperationalRate
	}
	return operationalRate
}

func (k Keeper) isAConsistentPowerOnPeriod(ctx sdk.Context, powerOnPeriod sdk.Dec) bool {
	smallestConsistentPowerOnPeriod := k.calculateSmallestConsistentPowerOnPeriod(ctx)
	// if powerOnPeriod >= smallestConsistentPowerOnPeriod
	if powerOnPeriod.GTE(smallestConsistentPowerOnPeriod) {
		return true
	}
	return false
}

func (k Keeper) calculateSmallestConsistentPowerOnPeriod(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)
	minimumPowerOnPeriodDec := sdk.NewDec(int64(params.MinimumPowerOnPeriod))
	maximumPowerOnPeriodDec := sdk.NewDec(int64(params.MaximumPowerOnPeriod))

	// smallestConsistentPowerOnPeriod = minimumPowerOnPeriodDec/maximumPowerOnPeriodDec
	smallestConsistentPowerOnPeriod := minimumPowerOnPeriodDec.Quo(maximumPowerOnPeriodDec)

	return smallestConsistentPowerOnPeriod
}
