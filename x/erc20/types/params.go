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
)

// Parameter store key
var (
	ParamStoreKeyEnableErc20   = []byte("EnableErc20")
	ParamStoreKeyEnableEVMHook = []byte("EnableEVMHook")
)

// NewParams creates a new Params object
func NewParams(
	enableErc20 bool,
	enableEVMHook bool,
) Params {
	return Params{
		EnableErc20:   enableErc20,
		EnableEVMHook: enableEVMHook,
	}
}

func DefaultParams() Params {
	return Params{
		EnableErc20:   true,
		EnableEVMHook: true,
	}
}

func ValidateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func (p Params) Validate() error {
	if err := ValidateBool(p.EnableEVMHook); err != nil {
		return err
	}

	return ValidateBool(p.EnableErc20)
}
