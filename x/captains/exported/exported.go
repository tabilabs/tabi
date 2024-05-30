package exported

import (
	"github.com/tabilabs/tabi/x/captains/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Captains interface {
	GetParams(ctx sdk.Context) types.Params

	// GetNodesByOwner returns all nodes owned by an address.
	GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []types.Node)

	// GetCurrentEpoch return the current epoch id.
	GetCurrentEpoch(ctx sdk.Context) uint64

	// CalcNodeCumulativeEmissionByEpoch returns the historical mission of the node at the end of epoch.
	CalcNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec

	// GetNodeClaimedEmission returns the historical emission user claimed.
	GetNodeClaimedEmission(ctx sdk.Context, nodeID string) sdk.Dec

	// UpdateGlobalAndNodeClaimedEmission updates the node_historical_emission_on_last_claim.
	UpdateGlobalAndNodeClaimedEmission(ctx sdk.Context, nodeID string) error
}
