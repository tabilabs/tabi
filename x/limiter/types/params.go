package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	// KeyEnabled is store's key for SendEnabled Params
	KeyEnabled = []byte("Enabled")
	// KeyWhiteList is store's key for the DefaultSendEnabled option
	KeyWhiteList = []byte("WhiteList")
)

// ParamKeyTable for limiter module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams returns a new Params object
func NewParams(enabled bool, whiteList []string) Params {
	return Params{
		Enabled:   enabled,
		WhiteList: whiteList,
	}
}

// DefaultParams is the default parameter configuration for the bank module
func DefaultParams() *Params {
	return &Params{
		Enabled:   false,
		WhiteList: []string{},
	}
}

// ParamSetPairs for limiter module
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnabled, &p.Enabled, validateIsBool),
		paramtypes.NewParamSetPair(KeyWhiteList, &p.WhiteList, validateWhiteList),
	}
}

// validateWhiteList validates a list of addresses
func validateWhiteList(i any) error {
	list, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	seenMap := make(map[string]bool)
	for _, addr := range list {
		if _, ok := seenMap[addr]; ok {
			return fmt.Errorf("duplicate whitelist address on %s", addr)
		}
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			return err
		}
		seenMap[addr] = true
	}
	return nil
}

// validateIsBool validates if the parameter is a boolean
func validateIsBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
