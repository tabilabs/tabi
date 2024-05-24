package keeper_test

import (
	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *IntegrationTestSuite) TestSaveDivision() {
	template := types.Division{
		Id:                       "division-template",
		Level:                    1,
		InitialSupply:            10000,
		SoldCount:                0,
		TotalCount:               0,
		ComputingPowerLowerBound: 10000,
		ComputingPowerUpperBound: 19999,
	}

	testCases := []struct {
		name      string
		division  types.Division
		execute   func(division types.Division) error
		expectErr bool
	}{
		{
			name:     "success - save division",
			division: template,
			execute: func(division types.Division) error {
				err := suite.Keeper.SaveDivision(suite.Ctx, division)
				suite.Require().NoError(err)

				return err
			},
			expectErr: false,
		},
		{
			name: "fail - save division duplicated",
			execute: func(division types.Division) error {
				err := suite.Keeper.SaveDivision(suite.Ctx, division)
				suite.Require().NoError(err)

				err = suite.Keeper.SaveDivision(suite.Ctx, division)
				suite.Require().Error(err)

				return err
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.division.Id = tc.name
			err := tc.execute(tc.division)
			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}

			res, found := suite.Keeper.GetDivision(suite.Ctx, tc.division.Id)
			suite.Require().True(found)
			suite.Require().Equal(tc.division, res)
		})
	}
}

func (suite *IntegrationTestSuite) TestDecideDivision() {
	testCases := []struct {
		name        string
		power       uint64
		expectLevel uint64
		expectErr   bool
	}{
		{
			name:        "fail - no division",
			power:       1000,
			expectLevel: 0,
		},
		{
			name:        "success - level 1",
			power:       9999,
			expectLevel: 1,
		},
		{
			name:        "success - level 2",
			power:       34999,
			expectLevel: 2,
		},
		{
			name:        "success - level 3",
			power:       104999,
			expectLevel: 3,
		},
		{
			name:        "success - level 4",
			power:       629999,
			expectLevel: 4,
		},
		{
			name:        "success - level 5",
			power:       100000000,
			expectLevel: 5,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			division := suite.Keeper.DecideDivision(suite.Ctx, tc.power)
			suite.Require().Equal(tc.expectLevel, division.Level)
		})
	}
}
