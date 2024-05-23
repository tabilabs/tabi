package keeper_test

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

const (
	EpochPhaseBeginEpoch = iota
	EpochPhaseStandBy
	EpochPhaseBeforeDigest
	EpochPhaseOnDigest
	EpochPhaseAfterDigest
	EpochPhaseBeforeBatches
	EpochPhaseOnBatches
	EpochPhaseAfterBatches
	EpochPhaseOnEnd
	EpochPhaseUpperBound
)

// EpochTestCase represents a test case for the epoch state.
type EpochTestCase struct {
	name     string // test name
	maxEpoch int    // number of epoches to run

	currEpochState *EpochState
	reporter       *CaptainsReporter

	execStandByFn func(*EpochState) // test case specific stand by function

	saveState bool                   // save state for each epoch-phase
	stateMap  map[string]*EpochState // state map
}

func (etc *EpochTestCase) Execute() {
	// execute as per epoch phase
	for etc.currEpochState.Epoch <= uint64(etc.maxEpoch) {
		// execute by phase
		etc.execute()
		// save state if needed
		if etc.saveState {
			etc.stateMap[etc.StateKey(etc.currEpochState.Epoch, etc.currEpochState.Phase)] = etc.currEpochState
		}
		etc.currEpochState.suite.T().Logf("execute epoch %d on phase %d", etc.currEpochState.Epoch, etc.currEpochState.Phase)
		// commit height
		etc.currEpochState.suite.Commit()
		// prepare the epoch state for the next phase
		etc.currEpochState = etc.Duplicate(etc.currEpochState)
		etc.currEpochState.Phase = (etc.currEpochState.Phase + 1) % EpochPhaseUpperBound
	}
}

func (etc *EpochTestCase) execute() {
	switch etc.currEpochState.Phase {
	case EpochPhaseBeginEpoch:
		execEpochPhaseBeginEpoch(etc.currEpochState)
	case EpochPhaseStandBy:
		if etc.execStandByFn != nil {
			etc.execStandByFn(etc.currEpochState)
		}
	case EpochPhaseBeforeDigest:
		execEpochPhaseBeforeDigest(etc.currEpochState)
	case EpochPhaseOnDigest:
		execEpochPhaseOnDigest(etc.currEpochState, etc.reporter)
	case EpochPhaseAfterDigest:
		execEpochPhaseAfterDigest(etc.currEpochState)
	case EpochPhaseBeforeBatches:
		execEpochPhaseBeforeBatches(etc.currEpochState)
	case EpochPhaseOnBatches:
		execEpochPhaseOnBatches(etc.currEpochState, etc.reporter)
	case EpochPhaseAfterBatches:
		execEpochPhaseAfterBatches(etc.currEpochState)
	case EpochPhaseOnEnd:
		execEpochPhaseOnEnd(etc.currEpochState, etc.reporter)
	default:
		panic("unknown epoch phase")
	}
}

// StateKey returns the mapping key for the epoch state.
func (etc *EpochTestCase) StateKey(epoch, phase uint64) string {
	return fmt.Sprintf("%d-%d", epoch, phase)
}

// Duplicate duplicates the epoch state.
func (etc *EpochTestCase) Duplicate(es *EpochState) *EpochState {
	nes := NewEpochState(es.suite)

	nes.Owner = es.Owner
	nes.Epoch = es.Epoch
	nes.Phase = es.Phase
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

func (etc *EpochTestCase) Export() {
	// TODO: print states as we may want to export it to a csv file
}

func execEpochPhaseBeginEpoch(nes *EpochState) {
	// update epoch
	nes.Epoch = nes.suite.Keeper.GetCurrentEpoch(nes.suite.Ctx)
	// check stand by over flag
	found := nes.suite.Keeper.HasStandByOverFlag(nes.suite.Ctx)
	nes.suite.Require().Equal(false, found)

	if nes.Epoch <= 2 {
		return
	}

	// NOTE: when we are able to check, we have already entered the next epoch.
	// so, the expected epoch should be epoch-2.
	found = nes.suite.Keeper.HasEpochEmission(nes.suite.Ctx, nes.Epoch-2)
	nes.suite.Require().Equal(false, found)

	found = nes.suite.Keeper.HasGlobalPledge(nes.suite.Ctx, nes.Epoch-2)
	nes.suite.Require().Equal(false, found)

	// TODO: has report digest, batch, end
}

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
	nes.GlobalComputingPower = sdk.ZeroDec() // reset
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
	crp.SubmitBatches(nes.suite, nes)
}

func execEpochPhaseAfterBatches(nes *EpochState) {
	for _, n := range nes.Nodes {
		if nes.Epoch > 1 {
			// check node cumulative emission 1 epoch before
			cee := nes.suite.Keeper.GetNodeCumulativeEmissionByEpoch(nes.suite.Ctx, nes.Epoch-1, n.ID)
			nes.suite.Require().Equal(cee, nes.NodeCumulativeEmission[n.ID])
		}
		if nes.Epoch > 2 {
			// check node accumulative epoch 2 epoch before is deleted
			found := nes.suite.Keeper.HasNodeCumulativeEmissionByEpoch(nes.suite.Ctx, nes.Epoch-2, n.ID)
			nes.suite.Require().Equal(false, found)
		}
		// check node computing power
		found := nes.suite.Keeper.HasNodeComputingPowerOnEpoch(nes.suite.Ctx, nes.Epoch, n.ID)
		nes.suite.Require().Equal(true, found)

		ncp := nes.suite.Keeper.GetNodeComputingPowerOnEpoch(nes.suite.Ctx, nes.Epoch, n.ID)
		nes.suite.Require().Equal(nes.NodeComputingPower[n.ID], ncp)
	}
	// check global computing power
	scp := nes.suite.Keeper.GetGlobalComputingPowerOnEpoch(nes.suite.Ctx, nes.Epoch)
	nes.suite.Require().Equal(nes.GlobalComputingPower, scp)

	if nes.Epoch > 2 {
		found := nes.suite.Keeper.HasOwnerPledge(nes.suite.Ctx, sdk.MustAccAddressFromBech32(nes.Owner), nes.Epoch-2)
		nes.suite.Require().Equal(false, found)
	}
}

func execEpochPhaseOnEnd(nes *EpochState, crp *CaptainsReporter) {
	crp.SubmitEnd(nes.suite, nes)
}

// Node represents a node in the captains module.
type Node struct {
	ID           string
	PowerOnRatio sdk.Dec
}

// Nodes represents a list of nodes.
type Nodes []Node

func (nds Nodes) PowerOnRatios(start, end int) []types.NodePowerOnRatio {
	ratios := make([]types.NodePowerOnRatio, end-start)
	for i := start; i < end; i++ {
		ratios[i-start] = types.NodePowerOnRatio{
			NodeId:           nds[i].ID,
			OnOperationRatio: nds[i].PowerOnRatio,
		}
	}
	return ratios
}

// EpochState captures the state of the captains module on specific height and epoch.
// NOTE: currently we set only one owner for all the nodes.
type EpochState struct {
	suite *IntegrationTestSuite

	// epoch phase
	Epoch uint64
	Phase uint64
	Owner string

	// epoch flag
	EndOnFlag bool

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
		suite:                   suite,
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

// WithNodes initializes the nodes in the epoch state.
func (es *EpochState) WithNodes(owner string, divisionLevel, amount uint64) *EpochState {
	nodeIds := es.suite.utilsBatchCreateCaptainNode(owner, divisionLevel, amount)

	nodes := make([]Node, len(nodeIds))
	for i, id := range nodeIds {
		nodes[i] = Node{
			ID:           id,
			PowerOnRatio: sdk.ZeroDec(),
		}
	}
	es.Nodes = nodes
	es.Owner = owner

	return es
}

// WithNodesPowerOnRatio initializes the power on ratio for the nodes in the epoch state.
func (es *EpochState) WithNodesPowerOnRatio() *EpochState {
	for i := range es.Nodes {
		// TODO: though we want an zero ratio as well.
		es.Nodes[i].PowerOnRatio = sdk.MustNewDecFromStr(fmt.Sprintf("%f", 0.47+rand.Float64()*0.53))
	}
	return es
}
