package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

type Captains interface {
	// GetNodesByOwner returns all nodes owned by an address.
	GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []types.Node)

	// GetCurrentEpoch return the current epoch id.
	GetCurrentEpoch(ctx sdk.Context) uint64

	// GetNodeHistoricalEmissionOnEpoch returns the historical mission of the node at the end of epoch.
	// NOTE: in epoch(t), we can get at historical emission by then end of epoch(t-2).
	GetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec

	// GetNodeHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
	// That says, if the user claims on epoch(T+1), the user will get rewards accrued at the end of epoch(T).
	// The next time the user claims, let's say on epoch(T+k+1), we can easily calc the rewards by subtracting
	// node_historical_emission_on_last_claim from node_historical_emission(T+k)
	GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec

	// UpdateNodeHistoricalEmissionOnLastClaim updates the node_historical_emission_on_last_claim.
	// NOTE: call this only after claiming.
	UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error
}
