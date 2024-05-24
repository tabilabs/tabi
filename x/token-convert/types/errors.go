package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrInvalidStrategy     = sdkerrors.Register(ModuleName, 1, "invalid strategy")
	ErrInvalidVoucher      = sdkerrors.Register(ModuleName, 2, "invalid voucher")
	ErrInvalidVoucherOwner = sdkerrors.Register(ModuleName, 3, "invalid voucher owner")
	ErrInsufficientFunds   = sdkerrors.Register(ModuleName, 4, "insufficient funds")
	ErrInvalidCoin         = sdkerrors.Register(ModuleName, 5, "invalid coin")
)
