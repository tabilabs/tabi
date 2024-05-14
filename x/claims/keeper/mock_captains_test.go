package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	captainnodetypes "github.com/tabilabs/tabi/x/captains/types"
)

type MockCaptains struct {
	ownerMap map[string]sdk.AccAddress
	nodeMap  map[string][]captainnodetypes.Node

	epoch uint64
}

func NewMockCaptains() *MockCaptains {
	return &MockCaptains{}
}

func (mock *MockCaptains) GetParams(ctx sdk.Context) captainnodetypes.Params {
	return captainnodetypes.DefaultParams()
}

func (mock *MockCaptains) GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []captainnodetypes.Node) {
	switch owner.String() {
	case Account01:
		return mock.nodeMap[Account01]
	case Account02:
		return mock.nodeMap[Account02]
	case Account03:
		return mock.nodeMap[Account03]
	}
	return []captainnodetypes.Node{}
}

// GetCurrentEpoch return the current epoch id.
func (mock *MockCaptains) GetCurrentEpoch(ctx sdk.Context) uint64 {
	return mock.epoch
}

func (mock *MockCaptains) SetCurrentEpoch(epoch uint64) {
	mock.epoch = epoch
}

// CalAndGetNodeHistoricalEmissionOnEpoch returns the historical mission of the node at the end of epoch.
func (mock *MockCaptains) CalAndGetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	return sdk.ZeroDec()
}

// GetNodeHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
func (mock *MockCaptains) GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec {
	return sdk.ZeroDec()
}

// UpdateNodeHistoricalEmissionOnLastClaim updates the node_historical_emission_on_last_claim.
// NOTE: call this only after claiming.
func (mock *MockCaptains) UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error {
	return nil
}
