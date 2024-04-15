package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParamsTestSuite struct {
	suite.Suite
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestParamsValidation() {
	testCases := []struct {
		name   string
		params Params
		expErr bool
	}{
		{
			name:   "NewParamsWithValidParams",
			params: DefaultParams(),
			expErr: false,
		},
		{
			name:   "NewParamsWithInvalidTotalCountCaptains",
			params: NewParams(0, 24, 6, 300000, 1, 5),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidMaximumPowerOnPeriod",
			params: NewParams(100000, 25, 6, 300000, 1, 5),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidMinimumPowerOnPeriod",
			params: NewParams(100000, 24, 25, 300000, 1, 5),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidConstantA",
			params: NewParams(100000, 24, 6, 0, 1, 5),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidCurrentLevelForSale",
			params: NewParams(100000, 24, 6, 300000, 6, 5),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidMaximumNumberOfHoldings",
			params: NewParams(100000, 24, 6, 300000, 1, 100001),
			expErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.params.Validate()
		if tc.expErr {
			suite.Require().Error(err, tc.name)
		} else {
			suite.Require().NoError(err, tc.name)
		}
	}
}
