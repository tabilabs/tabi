package types

import (
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	LevelOne = iota + 1
	LevelTwo
	LevelThree
	LevelFour
	LevelFive
)

// NewGenesisState constructs a GenesisState
func NewGenesisState(
	params Params,
	baseState BaseState,
	divisions []Division,
	nodes []Node,
	epochesEmission []EpochEmission,
	nodesClaimedEmission []NodeClaimedEmission,
	nodesCumulativeEmission []NodeCumulativeEmission,
	globalsPledge []GlobalPledge,
	ownersPledge []OwnerPledge,
	ownersClaimableComputingPower []ClaimableComputingPower,
	globalsComputingPower []GlobalComputingPower,
	nodesComputingPower []NodesComputingPower,
	batches []BatchBase,
) *GenesisState {
	return &GenesisState{
		Params:                        params,
		BaseState:                     baseState,
		Divisions:                     divisions,
		Nodes:                         nodes,
		EpochesEmission:               epochesEmission,
		NodesClaimedEmission:          nodesClaimedEmission,
		NodesCumulativeEmission:       nodesCumulativeEmission,
		GlobalsPledge:                 globalsPledge,
		OwnersPledge:                  ownersPledge,
		OwnersClaimableComputingPower: ownersClaimableComputingPower,
		GlobalsComputingPower:         globalsComputingPower,
		NodesComputingPower:           nodesComputingPower,
		Batches:                       batches,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:    DefaultParams(),
		BaseState: DefaultBaseState(),
		Divisions: DefaultDivision(),
	}
}

// DefaultBaseState returns a default BaseState
func DefaultBaseState() BaseState {
	return BaseState{
		EpochId:               1,
		NextNodeSequence:      1,
		GlobalClaimedEmission: sdk.ZeroDec(),
		IsStandBy:             true, // for sure we start from stand-by phase.
	}
}

// DefaultDivision returns a default Division
func DefaultDivision() []Division {
	return []Division{
		{
			Id:                       GenDivisionsId(LevelOne),
			Level:                    LevelOne,
			InitialSupply:            40_000,
			SoldCount:                0,
			TotalCount:               0,
			ComputingPowerLowerBound: 2_000,
			ComputingPowerUpperBound: 9_999,
		},
		{
			Id:                       GenDivisionsId(LevelTwo),
			Level:                    LevelTwo,
			InitialSupply:            60_000,
			SoldCount:                0,
			TotalCount:               0,
			ComputingPowerLowerBound: 10_000,
			ComputingPowerUpperBound: 39_999,
		},
		{
			Id:                       GenDivisionsId(LevelThree),
			Level:                    LevelThree,
			InitialSupply:            70_000,
			SoldCount:                0,
			TotalCount:               0,
			ComputingPowerLowerBound: 40_000,
			ComputingPowerUpperBound: 159_999,
		},
		{
			Id:                       GenDivisionsId(LevelFour),
			Level:                    LevelFour,
			InitialSupply:            20_000,
			SoldCount:                0,
			TotalCount:               0,
			ComputingPowerLowerBound: 160_000,
			ComputingPowerUpperBound: 959_999,
		},
		{
			Id:                       GenDivisionsId(LevelFive),
			Level:                    LevelFive,
			InitialSupply:            10_000,
			SoldCount:                0,
			TotalCount:               0,
			ComputingPowerLowerBound: 960_000,
			ComputingPowerUpperBound: math.MaxUint64,
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any.
func (gs *GenesisState) Validate() error {
	err := gs.Params.Validate()
	if err != nil {
		return err
	}

	err = gs.ValidateBaseState()
	if err != nil {
		return err
	}

	divIdMap, err := gs.ValidateDivisions()
	if err != nil {
		return err
	}

	nodesMap, err := gs.ValidateNodes(divIdMap)
	if err != nil {
		return err
	}

	err = gs.ValidateEpochEmission()
	if err != nil {
		return err
	}

	err = gs.ValidateNodesClaimedEmission(nodesMap)
	if err != nil {
		return err
	}

	err = gs.ValidateNodesCumulativeEmission(nodesMap)
	if err != nil {
		return err
	}

	err = gs.ValidateGlobalsPledge()
	if err != nil {
		return err
	}

	err = gs.ValidateOwnersPledge()
	if err != nil {
		return err
	}

	err = gs.ValidateOwnersClaimableComputingPower()
	if err != nil {
		return err

	}

	err = gs.ValidateGlobalsComputingPower()
	if err != nil {
		return err
	}

	err = gs.ValidateNodesComputingPower()
	if err != nil {
		return err
	}

	err = gs.ValidateBatches()
	if err != nil {
		return err
	}

	return nil
}

// ValidateBaseState performs basic base state validation returning an error upon any.
func (gs *GenesisState) ValidateBaseState() error {
	if gs.BaseState.EpochId == 0 {
		return fmt.Errorf("epoch id should be greter than zero, is %d", gs.BaseState.EpochId)
	}
	if gs.BaseState.NextNodeSequence == 0 {
		return fmt.Errorf("next node sequence should be greter than zero, is %d", gs.BaseState.NextNodeSequence)
	}
	return nil
}

// ValidateDivisions performs basic divisions validation returning an error upon any.
func (gs *GenesisState) ValidateDivisions() (map[string]uint64, error) {
	divMap := make(map[uint64]int)
	divIdMap := make(map[string]uint64)
	maxLevel := uint64(0)

	// check division level is unique
	for index, division := range gs.Divisions {
		if _, ok := divMap[division.Level]; ok {
			return nil, fmt.Errorf("duplicate division level %d", division.Level)
		}
		if division.Level == 0 {
			return nil, fmt.Errorf("division level should be greater than zero, is %d", division.Level)
		}
		if division.Id == "" {
			return nil, fmt.Errorf("division id is empty, in level %d", division.Level)
		}
		if division.InitialSupply == 0 {
			return nil, fmt.Errorf("inital supply should be greater than zero, is %d", division.InitialSupply)
		}
		if division.ComputingPowerLowerBound == 0 {
			return nil, fmt.Errorf("computing power lower bound should be greater than zero, is %d", division.ComputingPowerLowerBound)
		}
		if division.ComputingPowerUpperBound == 0 {
			return nil, fmt.Errorf("computing power upper bound should be greater than zero, is %d", division.ComputingPowerUpperBound)
		}
		if division.ComputingPowerLowerBound >= division.ComputingPowerUpperBound {
			return nil, fmt.Errorf("computing power lower bound should be less than upper bound, is %d", division.ComputingPowerLowerBound)
		}
		divMap[division.Level] = index
		divIdMap[division.Id] = division.Level
		if division.Level > maxLevel {
			maxLevel = division.Level
		}
	}
	// check division level is continuous
	if len(divMap) != int(maxLevel) {
		return nil, fmt.Errorf("division level is not continuous")
	}
	// check computing power bound is continuous
	for i := uint64(1); i < maxLevel; i++ {
		if gs.Divisions[divMap[i]].ComputingPowerUpperBound+1 != gs.Divisions[divMap[i+1]].ComputingPowerLowerBound {
			return nil, fmt.Errorf("computing power bound is not continuous, between level %d and %d", i, i+1)
		}
	}
	return divIdMap, nil
}

// ValidateNodes performs basic nodes validation returning an error upon any.
func (gs *GenesisState) ValidateNodes(divIdMap map[string]uint64) (map[string]bool, error) {
	seenMap := make(map[string]bool)
	for _, node := range gs.Nodes {
		if _, ok := seenMap[node.Id]; ok {
			return nil, fmt.Errorf("duplicate node id %s", node.Id)
		}
		if node.Id == "" {
			return nil, fmt.Errorf("node id is empty")
		}
		if node.Owner == "" {
			return nil, fmt.Errorf("node owner is empty")
		}
		if node.DivisionId == "" {
			return nil, fmt.Errorf("node division id is empty")
		}
		if node.ComputingPower == 0 {
			return nil, fmt.Errorf("node computing power should be greater than zero, is %d", node.ComputingPower)
		}
		if _, ok := divIdMap[node.DivisionId]; !ok {
			return nil, fmt.Errorf("unknown division id %s for node %s", node.DivisionId, node.Id)
		}
		seenMap[node.Id] = true
	}
	return seenMap, nil
}

// ValidateEpochEmission performs basic epoch emission validation returning an error upon any.
func (gs *GenesisState) ValidateEpochEmission() error {
	seenMap := make(map[uint64]bool)
	for _, ee := range gs.EpochesEmission {
		if _, ok := seenMap[ee.EpochId]; ok {
			return fmt.Errorf("duplicate epoch id %d", ee.EpochId)
		}
		if ee.EpochId == 0 {
			return fmt.Errorf("epoch id should be greater than zero, is %d", ee.EpochId)
		}
		seenMap[ee.EpochId] = true
	}
	return nil
}

// ValidateNodesClaimedEmission performs basic nodes claimed emission validation returning an error upon any.
func (gs *GenesisState) ValidateNodesClaimedEmission(nodesMap map[string]bool) error {
	seenMap := make(map[string]bool)
	for _, nce := range gs.NodesClaimedEmission {
		if !nodesMap[nce.NodeId] {
			return fmt.Errorf("unknown node id %s", nce.NodeId)
		}
		if _, ok := seenMap[nce.NodeId]; ok {
			return fmt.Errorf("duplicate node id %s", nce.NodeId)
		}
		if nce.Emission.IsZero() {
			return fmt.Errorf("emission is zero")
		}
		seenMap[nce.NodeId] = true
	}
	return nil
}

// ValidateNodesCumulativeEmission performs basic nodes cumulative emission validation returning an error upon any.
func (gs *GenesisState) ValidateNodesCumulativeEmission(nodesMap map[string]bool) error {
	seenMap := make(map[string]bool)
	for _, nce := range gs.NodesCumulativeEmission {
		if !nodesMap[nce.NodeId] {
			return fmt.Errorf("unknown node id %s", nce.NodeId)
		}
		if nce.EpochId > gs.BaseState.EpochId {
			return fmt.Errorf("epoch id %d is greater than current epoch id %d", nce.EpochId, gs.BaseState.EpochId)
		}
		if nce.EpochId == 0 {
			return fmt.Errorf("epoch id should be greater than zero, is %d", nce.EpochId)
		}
		if nce.NodeId == "" {
			return fmt.Errorf("node id is empty")
		}
		if nce.Emission.IsZero() {
			return fmt.Errorf("emission is zero")
		}
		uid := nce.NodeId + "-" + strconv.FormatUint(nce.EpochId, 10)
		if _, ok := seenMap[uid]; ok {
			return fmt.Errorf("duplicate on node id %s with epoch id %d", nce.NodeId, nce.EpochId)
		}
		seenMap[uid] = true
	}
	return nil
}

// ValidateGlobalsPledge performs basic globals pledge validation returning an error upon any.
// NOTE: no need to validate pledge amount as it can be zero.
func (gs *GenesisState) ValidateGlobalsPledge() error {
	seenMap := make(map[uint64]bool)
	for _, gp := range gs.GlobalsPledge {
		if _, ok := seenMap[gp.EpochId]; ok {
			return fmt.Errorf("duplicate epoch id %d", gp.EpochId)
		}
		if gp.EpochId == 0 {
			return fmt.Errorf("epoch id should be greater than zero, is %d", gp.EpochId)
		}
		seenMap[gp.EpochId] = true
	}
	return nil
}

// ValidateOwnersPledge performs basic owners pledge validation returning an error upon any.
func (gs *GenesisState) ValidateOwnersPledge() error {
	seenMap := make(map[string]bool)
	for _, op := range gs.OwnersPledge {
		if op.Owner == "" {
			return fmt.Errorf("owner is empty")
		}
		if op.EpochId == 0 {
			return fmt.Errorf("epoch id should be greater than zero, is %d", op.EpochId)
		}

		uid := op.Owner + "-" + strconv.FormatUint(op.EpochId, 10)
		if _, ok := seenMap[uid]; ok {
			return fmt.Errorf("duplicate owner %s", op.Owner)
		}
		seenMap[uid] = true
	}
	return nil
}

// ValidateOwnersClaimableComputingPower performs basic owners claimable computing power validation returning an error upon any.
func (gs *GenesisState) ValidateOwnersClaimableComputingPower() error {
	seenMap := make(map[string]bool)
	for _, occp := range gs.OwnersClaimableComputingPower {
		if occp.Owner == "" {
			return fmt.Errorf("owner is empty")
		}
		if _, ok := seenMap[occp.Owner]; ok {
			return fmt.Errorf("duplicate owner %s", occp.Owner)
		}
		seenMap[occp.Owner] = true
	}
	return nil
}

// ValidateGlobalsComputingPower performs basic globals computing power validation returning an error upon any.
func (gs *GenesisState) ValidateGlobalsComputingPower() error {
	seenMap := make(map[uint64]bool)
	for _, gcp := range gs.GlobalsComputingPower {
		if _, ok := seenMap[gcp.EpochId]; ok {
			return fmt.Errorf("duplicate epoch id %d", gcp.EpochId)
		}
		if gcp.EpochId == 0 {
			return fmt.Errorf("epoch id should be greater than zero, is %d", gcp.EpochId)
		}
		if gcp.Amount.IsZero() {
			return fmt.Errorf("computing power should be greater than zero, is %s", gcp.Amount)
		}
		seenMap[gcp.EpochId] = true
	}
	return nil
}

// ValidateNodesComputingPower performs basic nodes computing power validation returning an error upon any.
// NOTE: no need to validate computing power amount as it can be zero.
func (gs *GenesisState) ValidateNodesComputingPower() error {
	seenMap := make(map[string]bool)
	for _, ncp := range gs.NodesComputingPower {
		if ncp.NodeId == "" {
			return fmt.Errorf("node id is empty")
		}
		if ncp.EpochId == 0 {
			return fmt.Errorf("epoch id should be greater than zero, is %d", ncp.EpochId)
		}
		uid := ncp.NodeId + "-" + strconv.FormatUint(ncp.EpochId, 10)
		if _, ok := seenMap[uid]; ok {
			return fmt.Errorf("duplicate on node id %s with epoch id %d", ncp.NodeId, ncp.EpochId)
		}
		seenMap[uid] = true
	}
	return nil
}

// ValidateBatches performs basic batches validation returning an error upon any.
func (gs *GenesisState) ValidateBatches() error {
	seenMap := make(map[uint64]bool)
	for _, b := range gs.Batches {
		if b.BatchId == 0 {
			return fmt.Errorf("batch id should be greater than zero, is %d", b.BatchId)
		}
		if b.Count == 0 {
			return fmt.Errorf("batch count should be greater than zero, is %d", b.Count)
		}
		if !seenMap[b.BatchId] {
			return fmt.Errorf("duplicate batch id %d", b.BatchId)
		}
		seenMap[b.BatchId] = true
	}
	return nil
}
