package keeper_test

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

const (
	// TODO: some phrases can be merged so we can reduce deep copy times.
	EpochPhaseStandBy = iota + 1
	EpochPhaseBeforeDigest
	EpochPhaseOnDigest
	EpochPhaseAfterDigest
	EpochPhaseBeforeBatches
	EpochPhaseOnBatches
	EpochPhaseAfterBatches
	EpochPhaseOnEnd
	EpochPhaseNextEpoch
)

// EpochPhase represents the phase of the epoch.
type EpochPhase int

// Node represents a node in the captains module.
type Node struct {
	ID           string
	PowerOnRatio sdk.Dec
}

type Nodes []Node

// EpochState captures the state of the captains module on specific height and epoch.
// NOTE: currently we set only one owner for all the nodes.
type EpochState struct {
	suite *IntegrationTestSuite

	// epoch phase
	Epoch uint64
	Phase EpochPhase
	Owner string

	// epoch flag
	StandByOverFlag bool
	EndOnFlag       bool

	// nodes in the epoch
	Nodes Nodes

	// emission
	EpochEmission          sdk.Dec
	GlobalClaimedEmission  sdk.Dec
	NodeClaimedEmission    map[string]sdk.Dec
	NodeCumulativeEmission map[string]sdk.Dec

	// computing power
	ClaimableComputingPower map[string]sdk.Dec // owner -> power
	GlobalComputingPower    sdk.Dec
	NodeComputingPower      map[string]sdk.Dec

	// pledge
	GlobalPledge sdk.Dec
	OwnerPledge  sdk.Dec
}

func NewEpochState(suite *IntegrationTestSuite) *EpochState {
	return &EpochState{
		Epoch:                   suite.Keeper.GetCurrentEpoch(suite.Ctx),
		Phase:                   EpochPhaseStandBy,
		EpochEmission:           sdk.ZeroDec(),
		GlobalClaimedEmission:   sdk.ZeroDec(),
		NodeClaimedEmission:     make(map[string]sdk.Dec),
		NodeCumulativeEmission:  make(map[string]sdk.Dec),
		GlobalComputingPower:    sdk.ZeroDec(),
		ClaimableComputingPower: make(map[string]sdk.Dec),
		NodeComputingPower:      make(map[string]sdk.Dec),
		GlobalPledge:            sdk.ZeroDec(),
		OwnerPledge:             sdk.ZeroDec(),
	}
}

// Duplicate duplicates the epoch state.
func (es *EpochState) Duplicate() *EpochState {
	nes := NewEpochState(es.suite)

	nes.Owner = es.Owner
	nes.Epoch = es.Epoch
	nes.Phase = es.Phase
	nes.StandByOverFlag = es.StandByOverFlag
	nes.EndOnFlag = es.EndOnFlag

	nes.GlobalClaimedEmission = es.GlobalClaimedEmission
	nes.EpochEmission = es.EpochEmission
	nes.GlobalPledge = es.GlobalPledge
	nes.GlobalComputingPower = es.GlobalComputingPower
	nes.OwnerPledge = es.OwnerPledge

	nes.Nodes = make([]Node, len(es.Nodes))
	for i, n := range es.Nodes {
		nes.Nodes[i] = n
	}
	for k, v := range es.NodeClaimedEmission {
		nes.NodeClaimedEmission[k] = v
	}
	for k, v := range es.NodeCumulativeEmission {
		nes.NodeCumulativeEmission[k] = v
	}
	for k, v := range es.ClaimableComputingPower {
		nes.ClaimableComputingPower[k] = v
	}

	for k, v := range es.NodeComputingPower {
		nes.NodeComputingPower[k] = v
	}
	return nes
}

// Execute runs the full epoch state.
func ExecuteEpoch(oes *EpochState, cpr *CaptainsReporter) *EpochState {
	// deep copy the epoch state
	nes := oes.Duplicate()

	switch nes.Phase {
	case EpochPhaseStandBy:
		// note: claim reward, commit power, reward power etc.
		execEpochPhaseStandBy(nes)
	case EpochPhaseBeforeDigest:
		// note: check before digest
		execEpochPhaseBeforeDigest(nes)
	case EpochPhaseOnDigest:
		execEpochPhaseOnDigest(nes, cpr)
	case EpochPhaseAfterDigest:
		execEpochPhaseAfterDigest(nes)
	case EpochPhaseBeforeBatches:
		execEpochPhaseBeforeBatches(nes)
	case EpochPhaseOnBatches:
		execEpochPhaseOnBatches(nes, cpr)
	case EpochPhaseAfterBatches:
		execEpochPhaseAfterBatches(nes)
	case EpochPhaseOnEnd:
		execEpochPhaseOnEnd(nes, cpr)
	case EpochPhaseNextEpoch:
		// into next epoch
		nes.Epoch = nes.suite.Keeper.GetCurrentEpoch(nes.suite.Ctx)
		nes.suite.Require().Equal(oes.Epoch+1, nes.Epoch)
		execEpochNextEpoch(nes)
		nes.suite.Commit()
		nes.Phase = EpochPhaseStandBy
		return nes
	default:
		panic("unknown epoch phase")
	}

	// go to next height and phase
	nes.suite.Commit()
	nes.Phase++

	return nes
}

func execEpochPhaseStandBy(es *EpochState) {}

func execEpochPhaseBeforeDigest(nes *EpochState) {
	nes.EpochEmission = nes.suite.Keeper.CalcEpochEmission(nes.suite.Ctx, nes.Epoch, sdk.NewDecWithPrec(1, 0))

	if nes.Epoch <= 1 {
		return
	}

	// check
	exist := nes.suite.Keeper.HasGlobalPledge(nes.suite.Ctx, nes.Epoch)
	nes.suite.Require().Equal(true, exist)
	nes.GlobalPledge = nes.suite.Keeper.GetGlobalPledge(nes.suite.Ctx, nes.Epoch)
}

func execEpochPhaseOnDigest(nes *EpochState, crp *CaptainsReporter) {
	crp.SubmitDigest(nes.suite, nes)
}

func execEpochPhaseAfterDigest(nes *EpochState) {
	found := nes.suite.Keeper.HasStandByOverFlag(nes.suite.Ctx)
	nes.suite.Require().Equal(true, found)
	nes.StandByOverFlag = true

	ee := nes.suite.Keeper.GetEpochEmission(nes.suite.Ctx, nes.Epoch)
	nes.suite.Require().Equal(ee, nes.EpochEmission)

	if nes.Epoch <= 1 {
		return
	}

	exist := nes.suite.Keeper.HasGlobalPledge(nes.suite.Ctx, nes.Epoch)
	nes.suite.Require().Equal(false, exist)
	nes.GlobalPledge = sdk.ZeroDec() // reset
}

func execEpochPhaseBeforeBatches(nes *EpochState) {
	for _, n := range nes.Nodes {
		if nes.Epoch >= 2 {
			nes.NodeCumulativeEmission[n.ID] = nes.suite.Keeper.CalcNodeCumulativeEmissionByEpoch(nes.suite.Ctx, nes.Epoch-1, n.ID)

			found := nes.suite.Keeper.HasOwnerPledge(nes.suite.Ctx, sdk.MustAccAddressFromBech32(nes.Owner), nes.Epoch)
			nes.suite.Require().Equal(true, found)
			nes.OwnerPledge = nes.suite.Keeper.GetOwnerPledge(nes.suite.Ctx, sdk.MustAccAddressFromBech32(nes.Owner), nes.Epoch)
		}
		if nes.Epoch >= 3 {
			found := nes.suite.Keeper.HasNodeCumulativeEmissionByEpoch(nes.suite.Ctx, nes.Epoch-2, n.ID)
			nes.suite.Require().Equal(true, found)
		}

		nes.NodeComputingPower[n.ID] = nes.suite.Keeper.CalcNodeComputingPowerOnEpoch(nes.suite.Ctx, nes.Epoch, n.ID, n.PowerOnRatio)
		nes.GlobalComputingPower = nes.GlobalComputingPower.Add(nes.NodeComputingPower[n.ID])
	}
}

func execEpochPhaseOnBatches(nes *EpochState, crp *CaptainsReporter) {
	crp.SubmitDigest(nes.suite, nes)
}

func execEpochPhaseAfterBatches(nes *EpochState) {
	scp := sdk.NewDec(0)
	for _, n := range nes.Nodes {
		if nes.Epoch > 1 {
			cee := nes.suite.Keeper.GetNodeCumulativeEmissionByEpoch(nes.suite.Ctx, nes.Epoch-1, n.ID)
			nes.suite.Require().Equal(cee, nes.NodeCumulativeEmission[n.ID])
		}
		if nes.Epoch > 2 {
			found := nes.suite.Keeper.HasNodeCumulativeEmissionByEpoch(nes.suite.Ctx, nes.Epoch-2, n.ID)
			nes.suite.Require().Equal(false, found)
		}
		ncp := nes.suite.Keeper.GetNodeComputingPowerOnEpoch(nes.suite.Ctx, nes.Epoch, n.ID)
		nes.suite.Require().Equal(ncp, nes.NodeComputingPower[n.ID])
		scp = scp.Add(ncp)
	}

	nes.suite.Require().Equal(scp, nes.GlobalComputingPower)

	if nes.Epoch > 2 {
		found := nes.suite.Keeper.HasOwnerPledge(nes.suite.Ctx, sdk.MustAccAddressFromBech32(nes.Owner), nes.Epoch-1)
		nes.suite.Require().Equal(false, found)
	}
}

func execEpochPhaseOnEnd(nes *EpochState, crp *CaptainsReporter) {
	crp.SubmitEnd(nes.suite, nes)
}

func execEpochNextEpoch(nes *EpochState) {
	if nes.Epoch <= 1 {
		return
	}

	found := nes.suite.Keeper.HasEpochEmission(nes.suite.Ctx, nes.Epoch-1)
	nes.suite.Require().Equal(false, found)

	found = nes.suite.Keeper.HasGlobalPledge(nes.suite.Ctx, nes.Epoch-1)
	nes.suite.Require().Equal(false, found)

	found = nes.suite.Keeper.HasStandByOverFlag(nes.suite.Ctx)
	nes.suite.Require().Equal(false, found)

	// TODO: has report digest, batch, end
}

// InitNodes initializes the nodes in the epoch state.
func (es *EpochState) InitNodes(suite *IntegrationTestSuite, owner string, divisionLevel, amount uint64) {
	nodeIds := suite.utilsBatchCreateCaptainNode(owner, divisionLevel, amount)

	nodes := make([]Node, len(nodeIds))
	for i, id := range nodeIds {
		nodes[i] = Node{
			ID:           id,
			PowerOnRatio: sdk.ZeroDec(),
		}
	}
	es.Nodes = nodes
	es.Owner = owner
}

func (es *EpochState) InitNodesPowerOnRatio() {
	for i := range es.Nodes {
		// TODO: though we want an zero ratio as well.
		es.Nodes[i].PowerOnRatio = sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.47+rand.Float64()*0.53))
	}
}

func (nds Nodes) PowerOnRatios(start, end int) []types.NodePowerOnRatio {
	ratios := make([]types.NodePowerOnRatio, end-start)
	for i := start; i < end; i++ {
		ratios[i] = types.NodePowerOnRatio{
			NodeId:           nds[i].ID,
			OnOperationRatio: nds[i].PowerOnRatio,
		}
	}
	return ratios
}
