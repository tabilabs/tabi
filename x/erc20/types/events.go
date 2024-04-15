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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// erc20 events
const (
	EventTypeTokenLock             = "token_lock"
	EventTypeTokenUnlock           = "token_unlock"
	EventTypeMint                  = "mint"
	EventTypeConvertCoin           = "convert_coin"
	EventTypeConvertERC20          = "convert_erc20"
	EventTypeBurn                  = "burn"
	EventTypeRegisterCoin          = "register_coin"
	EventTypeRegisterERC20         = "register_erc20"
	EventTypeToggleTokenConversion = "toggle_token_conversion" // #nosec

	AttributeKeyCosmosCoin = "cosmos_coin"
	AttributeKeyERC20Token = "erc20_token" // #nosec
	AttributeKeyReceiver   = "receiver"

	// ERC20EventTransfer defines the transfer event for ERC20
	ERC20EventTransfer = "Transfer"
)

// LogTransfer Event type for Transfer(address from, address to, uint256 value)
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
