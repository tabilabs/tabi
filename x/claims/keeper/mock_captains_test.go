package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	captainnodetypes "github.com/tabilabs/tabi/x/captains/types"
)

type MockCaptains struct {
	caseNum int

	nodeHistoricalEmission            map[string]sdk.Dec
	NodeHistoricalEmissionOnLastClaim map[string]sdk.Dec
}

func NewMockCaptains(caseNum int) *MockCaptains {
	mock := &MockCaptains{
		caseNum: caseNum,
	}
	switch caseNum {
	case KeyCase06:
		mock.setUpForKeyCase06()
		return mock
	case KeyCase07:
		mock.setUpForKeyCase07()
		return mock
	case KeyCase08:
		mock.setUpForKeyCase08()
		return mock
	case KeyCase09:
		mock.setUpForKeyCase09()
		return mock
	case KeyQueryNodeTotalRewards03:
		mock.setUpForKeyQueryNodeTotalRewards03()
		return mock
	case KeyQueryNodeTotalRewards04:
		mock.setUpForKeyQueryNodeTotalRewards04()
		return mock
	case KeyQueryHolderTotalRewards04:
		mock.setUpForKeyQueryHolderTotalRewards04()
		return mock
	case KeyQueryHolderTotalRewards05:
		mock.setUpForKeyQueryHolderTotalRewards05()
		return mock
	case KeyQueryHolderTotalRewards06:
		mock.setUpForKeyQueryHolderTotalRewards06()
		return mock
	case KeyQueryHolderTotalRewards07:
		mock.setUpForKeyQueryHolderTotalRewards07()
		return mock
	default:
		return mock
	}
}

func (mock *MockCaptains) setUpForKeyCase06() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.ZeroDec(),
		"node2": sdk.ZeroDec(),
		"node3": sdk.ZeroDec(),
		"node4": sdk.ZeroDec(),
		"node5": sdk.ZeroDec(),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.ZeroDec(),
		"node2": sdk.ZeroDec(),
		"node3": sdk.ZeroDec(),
		"node4": sdk.ZeroDec(),
		"node5": sdk.ZeroDec(),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyCase07() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(100),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.ZeroDec(),
		"node2": sdk.ZeroDec(),
		"node3": sdk.ZeroDec(),
		"node4": sdk.ZeroDec(),
		"node5": sdk.ZeroDec(),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyCase08() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(100),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(100),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyCase09() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(200),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(200),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(200),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(0),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(0),
		"node5": sdk.NewDec(100),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyQueryNodeTotalRewards03() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyQueryNodeTotalRewards04() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(0),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

// HolderTotalRewards setUp start

func (mock *MockCaptains) setUpForKeyQueryHolderTotalRewards04() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyQueryHolderTotalRewards05() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(0),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyQueryHolderTotalRewards06() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(100),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(100),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

func (mock *MockCaptains) setUpForKeyQueryHolderTotalRewards07() {
	nodeHistoricalEmission := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(100),
	}

	nodeHistoricalEmissionOnLastClaim := map[string]sdk.Dec{
		"node1": sdk.NewDec(100),
		"node2": sdk.NewDec(100),
		"node3": sdk.NewDec(100),
		"node4": sdk.NewDec(100),
		"node5": sdk.NewDec(0),
	}
	mock.nodeHistoricalEmission = nodeHistoricalEmission
	mock.NodeHistoricalEmissionOnLastClaim = nodeHistoricalEmissionOnLastClaim
}

// HolderTotalRewards setUp end

func (mock *MockCaptains) GetParams(ctx sdk.Context) captainnodetypes.Params {
	return captainnodetypes.DefaultParams()
}

func (mock *MockCaptains) GetNodesByOwner(ctx sdk.Context, owner sdk.AccAddress) (nodes []captainnodetypes.Node) {
	switch mock.caseNum {
	case KeyCase01:
		return []captainnodetypes.Node{}
	case KeyCase02, KeyCase03, KeyCase04, KeyCase05, KeyQueryHolderTotalRewards04, KeyQueryHolderTotalRewards05:
		return []captainnodetypes.Node{
			{
				Id:             "node1",
				DivisionId:     "division1",
				Owner:          owner.String(),
				ComputingPower: 10000,
			},
		}
	case KeyCase06, KeyCase07, KeyCase08, KeyCase09, KeyQueryHolderTotalRewards06, KeyQueryHolderTotalRewards07:
		return []captainnodetypes.Node{
			{
				Id:             "node1",
				DivisionId:     "division1",
				Owner:          owner.String(),
				ComputingPower: 10000,
			},
			{
				Id:             "node2",
				DivisionId:     "division1",
				Owner:          owner.String(),
				ComputingPower: 10000,
			},
			{
				Id:             "node3",
				DivisionId:     "division1",
				Owner:          owner.String(),
				ComputingPower: 10000,
			},
			{
				Id:             "node4",
				DivisionId:     "division1",
				Owner:          owner.String(),
				ComputingPower: 10000,
			},
			{
				Id:             "node5",
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
	case KeyCase03, KeyCase04, KeyCase07, KeyCase08:
		return 2
	case KeyQueryNodeTotalRewards01,
		KeyQueryNodeTotalRewards03,
		KeyQueryNodeTotalRewards04:
		return 2
	case KeyQueryHolderTotalRewards01,
		KeyQueryHolderTotalRewards02,
		KeyQueryHolderTotalRewards03,
		KeyQueryHolderTotalRewards04,
		KeyQueryHolderTotalRewards05,
		KeyQueryHolderTotalRewards06,
		KeyQueryHolderTotalRewards07:
		return 2
	case KeyCase05, KeyCase09:
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
	case KeyCase06:
		return mock.nodeHistoricalEmission[nodeID]
	case KeyCase07:
		return mock.nodeHistoricalEmission[nodeID]
	case KeyCase08:
		return mock.nodeHistoricalEmission[nodeID]
	case KeyCase09:
		return mock.nodeHistoricalEmission[nodeID]
	default:
		return sdk.ZeroDec()
	}
}

// GetNodeClaimedEmission returns the historical emission the last time user claimed.
func (mock *MockCaptains) GetNodeClaimedEmission(ctx sdk.Context, nodeID string) sdk.Dec {
	switch mock.caseNum {
	case KeyCase02:
		return sdk.ZeroDec()
	case KeyCase03:
		return sdk.ZeroDec()
	case KeyCase04:
		return sdk.NewDec(100)
	case KeyCase05:
		return sdk.NewDec(100)
	case KeyCase06:
		return mock.NodeHistoricalEmissionOnLastClaim[nodeID]
	case KeyCase07:
		return mock.NodeHistoricalEmissionOnLastClaim[nodeID]
	case KeyCase08:
		return mock.NodeHistoricalEmissionOnLastClaim[nodeID]
	case KeyCase09:
		return mock.NodeHistoricalEmissionOnLastClaim[nodeID]
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
