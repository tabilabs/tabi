package types

import (
	"errors"
	fmt "fmt"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	evm "github.com/tabilabs/tabi/x/evm/types"

	"gopkg.in/yaml.v2"
)

// default paramspace for params keeper
const (
	DefaultParamSpace = "mint"
)

// Parameter store key
var (
	// params store for inflation params
	KeyInflation             = []byte("Inflation")
	KeyMintDenom             = []byte("MintDenom")
	KeyInflationDistribution = []byte("InflationDistribution")
)

var (
	DefaultMintDenom             = evm.DefaultEVMDenom
	DefaultInflation             = sdk.NewDecWithPrec(20, 2) //20%
	DefaultInflationDistribution = InflationDistribution{
		StakingRewards: sdk.NewDecWithPrec(2500000000, 10), // 25%
		ClaimsRewards:  sdk.NewDecWithPrec(7500000000, 10), // 75%
	}
)

// ParamTable for mint module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	mintDenom string,
	inflation sdk.Dec,
	inflationDistribution InflationDistribution,
) Params {
	return Params{
		MintDenom:             mintDenom,
		Inflation:             inflation,
		InflationDistribution: inflationDistribution,
	}
}

// DefaultParams returns default minting module parameters
func DefaultParams() Params {
	return Params{
		Inflation:             DefaultInflation,
		MintDenom:             DefaultMintDenom,
		InflationDistribution: DefaultInflationDistribution,
	}
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyInflation, &p.Inflation, validateInflation),
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyInflationDistribution, &p.InflationDistribution, validateInflationDistribution),
	}
}

// GetParamSpace implements params.ParamStruct
func (p *Params) GetParamSpace() string {
	return DefaultParamSpace
}

// Validate returns err if the Params is invalid
func (p Params) Validate() error {

	if err := validateInflation(p.Inflation); err != nil {
		return err
	}
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationDistribution(p.InflationDistribution); err != nil {
		return err
	}

	return nil
}

func validateInflation(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDecWithPrec(2, 1)) || v.LT(sdk.ZeroDec()) {
		return sdkerrors.Wrapf(
			ErrInvalidMintInflation,
			"Mint inflation [%s] should be between [0, 0.2] ",
			v.String(),
		)
	}

	return nil
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return sdkerrors.Wrapf(
			ErrInvalidMintDenom,
			"Mint denom [%s] should not be empty",
			v,
		)
	}

	return sdk.ValidateDenom(v)
}

func validateInflationDistribution(i interface{}) error {
	v, ok := i.(InflationDistribution)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.StakingRewards.IsNegative() {
		return errors.New("staking distribution ratio must not be negative")
	}

	if v.ClaimsRewards.IsNegative() {
		return errors.New("claims distribution ratio must not be negative")
	}

	return nil
}
