package types

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
			params: NewParams(0, 24, 6, sdk.NewDec(300000), sdk.NewDecWithPrec(16, 1), sdk.OneDec(), 5, nil),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidMaximumPowerOnPeriod",
			params: NewParams(100000, 25, 6, sdk.NewDec(300000), sdk.NewDecWithPrec(16, 1), sdk.OneDec(), 5, nil),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidMinimumPowerOnPeriod",
			params: NewParams(100000, 24, 25, sdk.NewDec(300000), sdk.NewDecWithPrec(16, 1), sdk.OneDec(), 5, nil),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidConstantA",
			params: NewParams(100000, 24, 6, sdk.NewDec(0), sdk.NewDecWithPrec(16, 1), sdk.OneDec(), 5, nil),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidCurrentLevelForSale",
			params: NewParams(100000, 24, 6, sdk.NewDec(300000), sdk.NewDecWithPrec(16, 1), sdk.OneDec(), 5, nil),
			expErr: true,
		},
		{
			name:   "NewParamsWithInvalidMaximumNumberOfHoldings",
			params: NewParams(100000, 24, 6, sdk.NewDec(300000), sdk.NewDecWithPrec(16, 1), sdk.OneDec(), 100001, nil),
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
