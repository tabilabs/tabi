package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	captainnodetypes "github.com/tabilabs/tabi/x/captains/types"
)

type MockCaptains struct {
	caseNum int

	historyOnLastClaim map[int]sdk.Dec
}

func NewMockCaptains(caseNum int) *MockCaptains {
	return &MockCaptains{caseNum: caseNum}
}

func (mock *MockCaptains) GetParams(ctx sdk.Context) captainnodetypes.Params {
	return captainnodetypes.DefaultParams()
}

func (mock *MockCaptains) GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []captainnodetypes.Node) {
	switch mock.caseNum {
	case KeyCase01:
		return []captainnodetypes.Node{}
	case KeyCase02, KeyCase03, KeyCase04, KeyCase05:
		return []captainnodetypes.Node{
			{
				Id:             "node1",
				DivisionId:     "division1",
				Owner:          owner.String(),
				ComputingPower: 10000,
			},
		}
	default:
		return []captainnodetypes.Node{}
	}

}

// GetCurrentEpoch return the current epoch id.
func (mock *MockCaptains) GetCurrentEpoch(ctx sdk.Context) uint64 {
	switch mock.caseNum {
	case KeyCase01:
		return 1
	case KeyCase02:
		return 1
	case KeyCase03:
		return 2
	case KeyCase04:
		return 2
	case KeyCase05:
		return 3
	default:
		return 1
	}
}

// CalAndGetNodeHistoricalEmissionOnEpoch returns the historical mission of the node at the end of epoch.
func (mock *MockCaptains) CalcAndGetNodeHistoricalEmissionOnEpoch(ctx sdk.Context, epochID uint64, nodeID string) sdk.Dec {
	switch mock.caseNum {
	case KeyCase02:
		return sdk.ZeroDec()
	case KeyCase03:
		return sdk.NewDec(100)
	case KeyCase04:
		return sdk.NewDec(100)
	case KeyCase05:
		return sdk.NewDec(200)
	default:
		return sdk.ZeroDec()
	}
}

// GetNodeHistoricalEmissionOnLastClaim returns the historical emission the last time user claimed.
func (mock *MockCaptains) GetNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) sdk.Dec {
	switch mock.caseNum {
	case KeyCase02:
		return sdk.ZeroDec()
	case KeyCase03:
		return sdk.ZeroDec()
	case KeyCase04:
		return sdk.NewDec(100)
	case KeyCase05:
		return sdk.NewDec(100)
	default:
		return sdk.ZeroDec()
	}
}

// UpdateNodeHistoricalEmissionOnLastClaim updates the node_historical_emission_on_last_claim.
// NOTE: call this only after claiming.
func (mock *MockCaptains) UpdateNodeHistoricalEmissionOnLastClaim(ctx sdk.Context, nodeID string) error {
	switch mock.caseNum {
	case KeyCase03:
		return nil
	case KeyCase05:
		return nil
	default:
		return nil
	}
}
