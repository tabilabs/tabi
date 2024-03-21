package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

// NewParams creates a new Params instance
func NewParams(
	totalCountCaptains uint64,
	maximumPowerOnPeriod uint64,
	minimumPowerOnPeriod uint64,
	constantA uint64,
	currentLevelForSale uint64,
	maximumNumberOfHoldings uint64,
) Params {
	return Params{
		TotalCountCaptains:      totalCountCaptains,
		MaximumPowerOnPeriod:    maximumPowerOnPeriod,
		MinimumPowerOnPeriod:    minimumPowerOnPeriod,
		ConstantA:               constantA,
		CurrentLevelForSale:     currentLevelForSale,
		MaximumNumberOfHoldings: maximumNumberOfHoldings,
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
		return fmt.Errorf("total count of captains should be positive")
	}

	if p.MaximumPowerOnPeriod <= 0 || p.MaximumPowerOnPeriod > 24 {
		return fmt.Errorf("maximum power on period should be positive and less than or equal to 24")
	}

	if p.MinimumPowerOnPeriod <= 0 || p.MaximumPowerOnPeriod > 24 || p.MinimumPowerOnPeriod > p.MaximumPowerOnPeriod {
		return fmt.Errorf("minimum power on period should be positive, less than or equal to 24 and less than or equal to maximum power on period")
	}

	if p.ConstantA <= 0 {
		return fmt.Errorf("constant A should be positive")
	}
	if p.CurrentLevelForSale <= 0 || p.CurrentLevelForSale > 5 {
		return fmt.Errorf("current level for sale should be positive and less than or equal to 7")
	}

	if p.MaximumNumberOfHoldings <= 0 || p.MaximumNumberOfHoldings > p.TotalCountCaptains {
		return fmt.Errorf("maximum number of holdings should be positive and less than or equal to total count of captains")
	}

	for _, caller := range p.Callers {
		if _, err := sdk.AccAddressFromBech32(caller); err != nil {
			return fmt.Errorf("caller address is invalid: %s", err)
		}
	}
	return nil
}
