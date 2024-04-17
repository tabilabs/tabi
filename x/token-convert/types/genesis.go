package types

// NewGenesisState create a module's genesis state.
func NewGenesisState(nextSeq uint64, strategies []Strategy, vouchers []Voucher) *GenesisState {
	return &GenesisState{
		VoucherSequence: nextSeq,
		Strategies:      nil,
		Vouchers:        nil,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(1, nil, nil)
}

// ValidateGenesis performs basic validation of genesis data returning an
func ValidateGenesis(data GenesisState) error {
	panic("impl me")
}
