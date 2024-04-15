package types

// NewGenesisState constructs a GenesisState
func NewGenesisState(params Params, fp FeePool) *GenesisState {
	return &GenesisState{
		Params:  params,
		FeePool: fp,
	}
}

// DefaultGenesisState gets raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		FeePool: InitialFeePool(),
	}
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(gs *GenesisState) error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}
	return gs.FeePool.ValidateGenesis()
}
