package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrDivisionExists              = errorsmod.Register(ModuleName, 2, "division already exists")
	ErrDivisionNotExists           = errorsmod.Register(ModuleName, 3, "division does not exist")
	ErrDivisionSoldOut             = errorsmod.Register(ModuleName, 4, "division is sold out")
	ErrInvalidDivisionID           = errorsmod.Register(ModuleName, 5, "invalid division id")
	ErrNodeExists                  = errorsmod.Register(ModuleName, 6, "node already exists")
	ErrUserHoldingQuantityExceeded = errorsmod.Register(ModuleName, 7, "user holding quantity exceeded")
	ErrUnauthorized                = errorsmod.Register(ModuleName, 8, "unauthorized address")
	ErrNodeNotExists               = errorsmod.Register(ModuleName, 9, "node does not exist")
	ErrInsufficientComputingPower  = errorsmod.Register(ModuleName, 10, "insufficient experience")
	ErrInvalidSaleLevel            = errorsmod.Register(ModuleName, 11, "new sale level must be greater than the current sale level")
	ErrInvalidCalculation          = errorsmod.Register(ModuleName, 12, "invalid calculation")
)
