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
		TotalCountCaptains:      100000,
		MaximumPowerOnPeriod:    24,
		MinimumPowerOnPeriod:    6,
		ConstantA:               300000,
		CurrentLevelForSale:     1,
		MaximumNumberOfHoldings: 5,
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
	if p.TotalCountCaptains <= 0 {
		return nil
	}

	if p.MaximumPowerOnPeriod <= 0 || p.MaximumPowerOnPeriod > 24 {
		return nil
	}

	if p.MinimumPowerOnPeriod <= 0 || p.MaximumPowerOnPeriod > 24 || p.MinimumPowerOnPeriod > p.MaximumPowerOnPeriod {
		return nil
	}

	if p.ConstantA <= 0 {
		return nil
	}
	if p.CurrentLevelForSale <= 0 || p.CurrentLevelForSale > 7 {
		return nil
	}
	if p.MaximumNumberOfHoldings <= 0 || p.MaximumNumberOfHoldings > p.TotalCountCaptains {
		return nil
	}
	return nil
}
