// SPDX-License-Identifier: Apache-2.0

// Package tests provides utility functions for testing.
package tests

import (
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

var (
	// UosmoDenomtrace represents the DenomTrace for the uosmo token.
	UosmoDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "uosmo",
	}
	// UosmoIbcdenom represents the IBC denomination for the uosmo token.
	UosmoIbcdenom = UosmoDenomtrace.IBCDenom()

	// UatomDenomtrace represents the DenomTrace for the uatom token.
	UatomDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-1",
		BaseDenom: "uatom",
	}
	// UatomIbcdenom represents the IBC denomination for the uatom token.
	UatomIbcdenom = UatomDenomtrace.IBCDenom()

	// UtabiDenomtrace represents the DenomTrace for the atabi token.
	UtabiDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "atabi",
	}
	// UtabiIbcdenom represents the IBC denomination for the atabi token.
	UtabiIbcdenom = UtabiDenomtrace.IBCDenom()

	// UatomOsmoDenomtrace represents the DenomTrace for the uatom token on the osmo channel.
	UatomOsmoDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0/transfer/channel-1",
		BaseDenom: "uatom",
	}
	// UatomOsmoIbcdenom represents the IBC denomination for the uatom token on the osmo channel.
	UatomOsmoIbcdenom = UatomOsmoDenomtrace.IBCDenom()

	// AtabiDenomtrace represents the DenomTrace for the atabi token.
	AtabiDenomtrace = transfertypes.DenomTrace{
		Path:      "transfer/channel-0",
		BaseDenom: "atabi",
	}
	// AtabiIbcdenom represents the IBC denomination for the atabi token.
	AtabiIbcdenom = AtabiDenomtrace.IBCDenom()
)

