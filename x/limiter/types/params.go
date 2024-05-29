package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewParams returns a new Params object
func NewParams(enabled bool, whiteList []string) Params {
	return Params{
		Enabled:   enabled,
		AllowList: whiteList,
	}
}

// DefaultParams is the default parameter configuration for the bank module
func DefaultParams() *Params {
	return &Params{
		Enabled:   false,
		AllowList: []string{},
	}
}

// ValidateParams validates the parameters
func ValidateParams(params *Params) error {
	if params == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "nil params")
	}

	seenMap := make(map[string]bool)
	for _, addr := range params.AllowList {
		if _, ok := seenMap[addr]; ok {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "duplicate whitelist address")
		}
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return errorsmod.Wrap(err, "invalid whitelist address")
		}
		seenMap[addr] = true
	}
	return nil
}
