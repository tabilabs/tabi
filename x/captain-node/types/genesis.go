package types

import "math"

const (
	LevelOne = iota + 1
	LevelTwo
	LevelThree
	LevelFour
	LevelFive
)

// NewGenesisState constructs a GenesisState
func NewGenesisState(params Params, divisions []*Division, entries []*Entry) *GenesisState {
	return &GenesisState{
		Params:    params,
		Divisions: divisions,
		Entries:   entries,
	}
}

// Validate validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func (gen GenesisState) Validate() error {
	if err := gen.Params.Validate(); err != nil {
		return err
	}
	return nil
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), DefaultDivision(), nil)
}

func DefaultDivision() []*Division {
	return newDivisions()
}

func newDivisions() []*Division {
	return []*Division{
		{
			Id:             GenDivisionsId(LevelOne),
			Level:          LevelOne,
			TotalCount:     20000,
			ComputingPower: 2000,
			Low:            0,
			High:           0,
		},
		{
			Id:             GenDivisionsId(LevelTwo),
			Level:          LevelTwo,
			TotalCount:     30000,
			ComputingPower: 10000,
			Low:            0,
			High:           0,
		},
		{
			Id:             GenDivisionsId(LevelThree),
			Level:          LevelThree,
			TotalCount:     35000,
			ComputingPower: 35000,
			Low:            0,
			High:           0,
		},
		{
			Id:             GenDivisionsId(LevelFour),
			Level:          LevelThree,
			TotalCount:     10000,
			ComputingPower: 105000,
			Low:            0,
			High:           0,
		},
		{
			Id:             GenDivisionsId(LevelFive),
			Level:          LevelFive,
			TotalCount:     5000,
			ComputingPower: 630000,
			Low:            0,
			High:           math.MaxUint64,
		},
	}
}
