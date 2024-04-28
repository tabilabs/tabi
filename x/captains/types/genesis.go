package types

import (
	"math"
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
	epochState EpochState,
	divisions []Division,
	nodes []Node,
	nodesExtraInfo []NodeExtraInfo,
	powers []ClaimableComputingPower,
) *GenesisState {
	return &GenesisState{
		Params:                   params,
		EpochState:               epochState,
		Divisions:                divisions,
		Nodes:                    nodes,
		NodesExtraInfo:           nodesExtraInfo,
		ClaimableComputingPowers: powers,
	}
}

// Validate validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func (gs *GenesisState) Validate() error {
	// FIXME: impl full validations
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		DefaultParams(),
		DefaultEpochState(),
		DefaultDivision(),
		nil,
		nil,
		nil)
}

func DefaultEpochState() EpochState {
	return EpochState{
		CurrEpoch: 1,
		IsEnd:     false,
		Digest:    nil,
		Batches:   nil,
		Current:   EpochBase{},
		Previous:  EpochBase{},
	}
}

func DefaultDivision() []Division {
	return newDivisions()
}

func newDivisions() []Division {
	return []Division{
		{
			Id:                       GenDivisionsId(LevelOne),
			Level:                    LevelOne,
			TotalCount:               20000,
			ComputingPowerLowerBound: 2000,
			ComputingPowerUpperBound: 9999,
		},
		{
			Id:                       GenDivisionsId(LevelTwo),
			Level:                    LevelTwo,
			TotalCount:               30000,
			ComputingPowerLowerBound: 10000,
			ComputingPowerUpperBound: 34999,
		},
		{
			Id:                       GenDivisionsId(LevelThree),
			Level:                    LevelThree,
			TotalCount:               35000,
			ComputingPowerLowerBound: 35000,
			ComputingPowerUpperBound: 104999,
		},
		{
			Id:                       GenDivisionsId(LevelFour),
			Level:                    LevelFour,
			TotalCount:               10000,
			ComputingPowerLowerBound: 105000,
			ComputingPowerUpperBound: 629999,
		},
		{
			Id:                       GenDivisionsId(LevelFive),
			Level:                    LevelFive,
			TotalCount:               5000,
			ComputingPowerLowerBound: 630000,
			ComputingPowerUpperBound: math.MaxUint64,
		},
	}
}
