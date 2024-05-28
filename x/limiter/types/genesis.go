package types

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
func ValidateGenesis(data GenesisState) error {
	return ValidateParams(data.Params)
}
