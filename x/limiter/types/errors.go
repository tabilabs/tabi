package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrMemberAlreadyExisted = errorsmod.Register(ModuleName, 2, "member already existed in allow list")
	ErrMemberNotFound       = errorsmod.Register(ModuleName, 3, "member not found in allow list")
	ErrEmptyAllowList       = errorsmod.Register(ModuleName, 4, "empty allow list")
)
