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

package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// prefix bytes for the inflation persistent store
const prefixEpochMintProvision = 2

// KeyPrefixEpochMintProvision key prefix
var KeyPrefixEpochMintProvision = []byte{prefixEpochMintProvision}

// MigrateStore migrates the x/inflation module state from the consensus version 2 to
// version 3. Specifically, it deletes the EpochMintProvision from the store
func MigrateStore(store sdk.KVStore) error {
	store.Delete(KeyPrefixEpochMintProvision)
	return nil
}
