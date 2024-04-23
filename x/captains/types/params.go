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

// ParamKeyTable for captains node
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		CaptainsTotalCount:    100000,
		MinimumPowerOnPeriod:  6,
		MaximumPowerOnPeriod:  24,
		CaptainsConstant:      300000,
		MaximumHoldingAmount:  5,
		HalvingEraCoefficient: 0,
		CurrentSaleLevel:      1,
		AuthorizedMembers:     nil,
	}
}

// NewParams creates a Params.
func NewParams(
	captainsTotalCount uint64,
	minimumPowerOnPeriod uint64,
	maximumPowerOnPeriod uint64,
	captainsConstant uint64,
	maximumHoldingAmount uint64,
	halvingEraCoefficient uint64,
	currentSaleLevel uint64,
	authorizedMembers []string,
) Params {
	return Params{
		CaptainsTotalCount:    captainsTotalCount,
		MinimumPowerOnPeriod:  minimumPowerOnPeriod,
		MaximumPowerOnPeriod:  maximumPowerOnPeriod,
		CaptainsConstant:      captainsConstant,
		MaximumHoldingAmount:  maximumHoldingAmount,
		HalvingEraCoefficient: halvingEraCoefficient,
		CurrentSaleLevel:      currentSaleLevel,
		AuthorizedMembers:     authorizedMembers,
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
func (p *Params) Validate() error {
	if p.CaptainsTotalCount <= 0 {
		return fmt.Errorf("total count of captains should be positive")
	}

	if p.MaximumPowerOnPeriod <= 0 || p.MaximumPowerOnPeriod > 24 {
		return fmt.Errorf("maximum power on period should be positive and less than or equal to 24")
	}

	if p.MinimumPowerOnPeriod <= 0 || p.MaximumPowerOnPeriod > 24 || p.MinimumPowerOnPeriod > p.MaximumPowerOnPeriod {
		return fmt.Errorf("minimum power on period should be positive, less than or equal to 24 and less than or equal to maximum power on period")
	}

	if p.CaptainsConstant <= 0 {
		return fmt.Errorf("captains constant should be positive")
	}
	if p.CurrentSaleLevel <= 0 || p.CurrentSaleLevel > 5 {
		return fmt.Errorf("current sale level should be non-negative and less than or equal to 7")
	}

	if p.MaximumHoldingAmount <= 0 || p.MaximumHoldingAmount > p.CaptainsTotalCount {
		return fmt.Errorf("maximum number of holdings should be positive and less than or equal to total count of captains")
	}

	for _, memeber := range p.AuthorizedMembers {
		if _, err := sdk.AccAddressFromBech32(memeber); err != nil {
			return fmt.Errorf("memeber address is invalid: %s", err)
		}
	}
	return nil
}
