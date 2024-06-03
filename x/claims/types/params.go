package types

import (
	tabitypes "github.com/tabilabs/tabi/types"
)

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		EnableClaims: true,
		ClaimsDenom:  tabitypes.AttoVeTabi,
	}
}

// Validate returns err if the Params is invalid
func (p Params) ValidateBasic() error {
	return nil
}
