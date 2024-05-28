package keeper_test

import (
	"github.com/tabilabs/tabi/x/limiter/types"
)

func (suite *IntegrationTestSuite) TestSetParams() {
	params := types.Params{
		Enabled:   true,
		AllowList: []string{accounts[0].String(), accounts[1].String()},
	}

	// test set and get
	suite.App.LimiterKeeper.SetParams(suite.Ctx, params)
	got := suite.App.LimiterKeeper.GetParams(suite.Ctx)
	suite.Require().Equal(params, got)

	// test query params
	resp, err := suite.QueryClient.Params(suite.Ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(params, resp.Params)

	// test is enabled
	suite.Require().Equal(params.Enabled, suite.App.LimiterKeeper.IsEnabled(suite.Ctx))

	// test is authorized
	suite.Require().Equal(true, suite.App.LimiterKeeper.IsAuthorized(suite.Ctx, accounts[0]))
	suite.Require().Equal(true, suite.App.LimiterKeeper.IsAuthorized(suite.Ctx, accounts[1]))
	suite.Require().Equal(false, suite.App.LimiterKeeper.IsAuthorized(suite.Ctx, accounts[2]))
}
