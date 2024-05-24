package keeper_test

import (
	"fmt"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestCalcNodeComputingPowerOnEpoch() {
	testCases := []struct {
		name         string
		pledge       sdk.Dec
		powerOnRatio sdk.Dec
	}{
		{
			name:         "case 1",
			pledge:       sdk.ZeroDec(),
			powerOnRatio: sdk.NewDecWithPrec(1, 0),
		},
		{
			name:         "case 2",
			pledge:       sdk.ZeroDec(),
			powerOnRatio: sdk.ZeroDec(),
		},
		{
			name:         "case 3",
			pledge:       sdk.MustNewDecFromStr("0.511111111"),
			powerOnRatio: sdk.MustNewDecFromStr("0.111111111"),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			suite.Require().NotPanics(
				func() {
					exponentiation := tc.pledge.Mul(sdk.NewDec(20)).Quo(sdk.NewDec(3)).MustFloat64()
					exponentiated := sdk.MustNewDecFromStr(fmt.Sprintf("%f", math.Exp(exponentiation)))
					sdk.NewDec(int64(2000)).Mul(exponentiated).Mul(tc.powerOnRatio)
				})
		})
	}
}
