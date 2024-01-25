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

import (
	fmt "fmt"

	epochstypes "github.com/tabilabs/tabi/v1/x/epochs/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	period uint64,
	epochIdentifier string,
	epochsPerPeriod int64,
	skippedEpochs uint64,
) GenesisState {
	return GenesisState{
		Params:          params,
		Period:          period,
		EpochIdentifier: epochIdentifier,
		EpochsPerPeriod: epochsPerPeriod,
		SkippedEpochs:   skippedEpochs,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:          DefaultParams(),
		Period:          uint64(0),
		EpochIdentifier: epochstypes.DayEpochID,
		EpochsPerPeriod: 365,
		SkippedEpochs:   0,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := epochstypes.ValidateEpochIdentifierInterface(gs.EpochIdentifier); err != nil {
		return err
	}

	if err := validateEpochsPerPeriod(gs.EpochsPerPeriod); err != nil {
		return err
	}

	if err := validateSkippedEpochs(gs.SkippedEpochs); err != nil {
		return err
	}

	return gs.Params.Validate()
}

func validateEpochsPerPeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid genesis state type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("epochs per period must be positive: %d", v)
	}

	return nil
}

func validateSkippedEpochs(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid genesis state type: %T", i)
	}
	return nil
}
