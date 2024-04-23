package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrCalculateRewards                        = errorsmod.Register(ModuleName, 2, "error while calculating rewards")
	ErrSendCoins                               = errorsmod.Register(ModuleName, 3, "error while sending coins")
	ErrUpdateNodeHistoricalEmissionOnLastClaim = errorsmod.Register(ModuleName, 4, "error while updating node historical emission on last claim")
	ErrHolderNotFound                          = errorsmod.Register(ModuleName, 5, "holder not found")
)
