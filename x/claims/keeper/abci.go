package keeper

import (
	"time"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tabilabs/tabi/x/claims/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock updates node block reward
func (k *Keeper) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
}
