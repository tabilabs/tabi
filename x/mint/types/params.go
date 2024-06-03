package types

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	evm "github.com/tabilabs/tabi/x/evm/types"
)

var (
	DefaultMintDenom = evm.DefaultEVMDenom
	DefaultInflation = sdk.NewDecWithPrec(20, 2) // 20%
)

func NewParams(
	mintDenom string,
	inflation sdk.Dec,
) Params {
	return Params{
		MintDenom: mintDenom,
		Inflation: inflation,
	}
}

// DefaultParams returns default minting module parameters
func DefaultParams() Params {
	return Params{
		Inflation: DefaultInflation,
		MintDenom: DefaultMintDenom,
	}
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate returns err if the Params is invalid
func (p Params) Validate() error {
	if err := validateInflation(p.Inflation); err != nil {
		return err
	}
	if err := validateMintDenom(p.MintDenom); err != nil {
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
