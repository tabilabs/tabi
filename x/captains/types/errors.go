package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrDivisionExists    = errorsmod.Register(ModuleName, 2, "division already exists")
	ErrDivisionNotExists = errorsmod.Register(ModuleName, 3, "division does not exist")
	ErrDivisionSoldOut   = errorsmod.Register(ModuleName, 4, "division is sold out")
	ErrNodeExists        = errorsmod.Register(ModuleName, 5, "node already exists")
	ErrUnauthorized      = errorsmod.Register(ModuleName, 6, "unauthorized address")

	ErrNodeNotExists              = errorsmod.Register(ModuleName, 7, "node does not exist")
	ErrInsufficientComputingPower = errorsmod.Register(ModuleName, 8, "insufficient experience")
	ErrInvalidSaleLevel           = errorsmod.Register(ModuleName, 9, "new sale level must be greater than the current sale level")
	ErrInvalidReport              = errorsmod.Register(ModuleName, 10, "invalid report")
	ErrDeleteLastMember           = errorsmod.Register(ModuleName, 11, "can not delete the last member")
	ErrTypeOverflow               = errorsmod.Register(ModuleName, 12, "power overflow")
)
