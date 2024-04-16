package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrStrategyNotFound      = sdkerrors.Register(ModuleName, 1, "strategy not found")
	ErrVoucherNotFound       = sdkerrors.Register(ModuleName, 2, "voucher not found")
	ErrInvalidVoucherOwner   = sdkerrors.Register(ModuleName, 3, "invalid voucher owner")
	ErrStrategyAlreadyExists = sdkerrors.Register(ModuleName, 4, "strategy already exists")
	ErrInsufficientFunds     = sdkerrors.Register(ModuleName, 5, "insufficient funds")
)
