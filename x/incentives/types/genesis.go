// Copyright 2024 Tabi Foundation
// This file is part of the Tabi Network packages.
//
// Tabi is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Tabi packages are distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.

package types

import "fmt"

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	params Params,
	incentives []Incentive,
	gasMeters []GasMeter,
) GenesisState {
	return GenesisState{
		Params:     params,
		Incentives: incentives,
		GasMeters:  gasMeters,
	}
}

// DefaultGenesisState sets default evm genesis state with empty accounts and
// default params and chain config values.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	seenContractIn := make(map[string]bool)
	for _, in := range gs.Incentives {
		// only one incentive per contract
		if seenContractIn[in.Contract] {
			return fmt.Errorf("contract duplicated on genesis '%s'", in.Contract)
		}

		if err := in.Validate(); err != nil {
			return err
		}

		seenContractIn[in.Contract] = true
	}

	seenGasMeters := make(map[string]bool)
	for _, gm := range gs.GasMeters {
		// only one gas meter per contract+participant combination
		if seenGasMeters[gm.Contract+gm.Participant] {
			return fmt.Errorf(
				"gas meter duplicated on genesis contract: '%s',  participant: '%s'",
				gm.Contract, gm.Participant,
			)
		}

		if err := gm.Validate(); err != nil {
			return err
		}

		seenGasMeters[gm.Contract+gm.Participant] = true
	}

	return gs.Params.Validate()
}
