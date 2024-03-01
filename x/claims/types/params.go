package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// default paramspace for params keeper
const (
	DefaultParamSpace = ModuleName
)

// ParamTable for mint module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		EnableClaims: true,
		ClaimsDenom:  "claims",
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return nil
}

// GetParamSpace implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

// Validate returns err if the Params is invalid
func (p Params) Validate() error {
	return nil
}
