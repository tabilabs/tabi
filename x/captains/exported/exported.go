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

	// CalcAndGetNodeCumulativeEmissionByEpoch returns the historical mission of the node at the end of epoch.
	CalcAndGetNodeCumulativeEmissionByEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec

	// GetNodeClaimedEmission returns the historical emission user claimed.
	GetNodeClaimedEmission(ctx sdk.Context, nodeID string) sdk.Dec

	// UpdateGlobalAndNodeClaimedEmission updates the node_historical_emission_on_last_claim.
	UpdateGlobalAndNodeClaimedEmission(ctx sdk.Context, nodeID string) error
}
