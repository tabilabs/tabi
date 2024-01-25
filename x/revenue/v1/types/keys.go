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

import sdk "github.com/cosmos/cosmos-sdk/types"

// constants
const (
	// module name
	ModuleName = "revenue"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the fees persistent store
const (
	prefixRevenue = iota + 1
	prefixDeployer
	prefixWithdrawer
)

// KVStore key prefixes
var (
	KeyPrefixRevenue    = []byte{prefixRevenue}
	KeyPrefixDeployer   = []byte{prefixDeployer}
	KeyPrefixWithdrawer = []byte{prefixWithdrawer}
)

// GetKeyPrefixDeployer returns the KVStore key prefix for storing
// registered revenue contract for a deployer
func GetKeyPrefixDeployer(deployerAddress sdk.AccAddress) []byte {
	return append(KeyPrefixDeployer, deployerAddress.Bytes()...)
}

// GetKeyPrefixWithdrawer returns the KVStore key prefix for storing
// registered revenue contract for a withdrawer
func GetKeyPrefixWithdrawer(withdrawerAddress sdk.AccAddress) []byte {
	return append(KeyPrefixWithdrawer, withdrawerAddress.Bytes()...)
}
