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
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrInternalRevenue              = errorsmod.Register(ModuleName, 2, "internal revenue error")
	ErrRevenueDisabled              = errorsmod.Register(ModuleName, 3, "revenue module is disabled by governance")
	ErrRevenueAlreadyRegistered     = errorsmod.Register(ModuleName, 4, "revenue already exists for given contract")
	ErrRevenueNoContractDeployed    = errorsmod.Register(ModuleName, 5, "no contract deployed")
	ErrRevenueContractNotRegistered = errorsmod.Register(ModuleName, 6, "no revenue registered for contract")
	ErrRevenueDeployerIsNotEOA      = errorsmod.Register(ModuleName, 7, "no revenue registered for contract")
)
