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

// constants
const (
	// module name
	ModuleName = "inflation"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the inflation persistent store
const (
	prefixPeriod = iota + 1
	prefixEpochMintProvision
	prefixEpochIdentifier
	prefixEpochsPerPeriod
	prefixSkippedEpochs
)

// KVStore key prefixes
var (
	KeyPrefixPeriod          = []byte{prefixPeriod}
	KeyPrefixEpochIdentifier = []byte{prefixEpochIdentifier}
	KeyPrefixEpochsPerPeriod = []byte{prefixEpochsPerPeriod}
	KeyPrefixSkippedEpochs   = []byte{prefixSkippedEpochs}
)
