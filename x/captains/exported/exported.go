package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

type Captains interface {
	GetParams(ctx sdk.Context) types.Params

	// GetNodesByOwner returns all nodes owned by an address.
	GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []types.Node)

	// GetCurrentEpoch return the current epoch id.
	GetCurrentEpoch(ctx sdk.Context) uint64

	// CalAndGetNodeHistoricalEmissionOnEpoch returns the historical mission of the node at the end of epoch.
	CalAndGetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec

	// GetNodeHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
	GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec

	// UpdateNodeHistoricalEmissionOnLastClaim updates the node_historical_emission_on_last_claim.
	// NOTE: call this only after claiming.
	UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error
}
