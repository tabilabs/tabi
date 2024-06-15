package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestTechCoefficient() {
	testCases := []struct {
		name   string
		level  uint64
		expect string
	}{
		{
			name:   "sale level 1",
			level:  1,
			expect: "1",
		},
		{
			name:   "sale level 2",
			level:  2,
			expect: "1.6",
		},
		{
			name:   "sale level 3",
			level:  3,
			expect: "2.56",
		},
		{
			name:   "sale level 4",
			level:  4,
			expect: "4.096",
		},
		{
			name:   "sale level 5",
			level:  5,
			expect: "6.5536",
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("TechCoefficient - %s", tc.name), func() {
			suite.utilsUpdateLevel(tc.level)

			tec := suite.Keeper.CalcTechProgressCoefficient(suite.Ctx)
			expect, err := sdk.NewDecFromStr(tc.expect)

			suite.Require().NoError(err)
			suite.Require().Equal(tec, expect)
		})
	}
}

func (suite *IntegrationTestSuite) TestEpochEmission() {
	testCases := []struct {
		name string
	}{
		{
			name: "test emission",
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("TechCoefficient - %s", tc.name), func() {
			res := suite.Keeper.CalcBaseEpochEmission(suite.Ctx).Quo(sdk.NewDec(1e18))
			suite.T().Log("BaseEmission: ", res.String())
			res = suite.Keeper.CalcEpochEmission(suite.Ctx, 1, sdk.MustNewDecFromStr("0.47")).Quo(sdk.NewDec(1e18))
			suite.T().Log("Emission: ", res.String())
		})
	}
}
