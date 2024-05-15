package keeper_test

import (
	"fmt"

	"github.com/tabilabs/tabi/x/claims/types"
)

func (suite *ClaimsTestSuite) TestQueryParams() {
	testCases := []struct {
		name    string
		expPass bool
	}{
		{
			"pass",
			true,
		},
	}
	for _, tc := range testCases {
		params := suite.app.ClaimsKeeper.GetParams(suite.ctx)
		exp := &types.QueryParamsResponse{Params: params}

		res, err := suite.queryClient.Params(suite.ctx.Context(), &types.QueryParamsRequest{})
		if tc.expPass {
			suite.Require().Equal(exp, res, tc.name)
			suite.Require().NoError(err)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *ClaimsTestSuite) TestQueryNodeTotalRewards() {
	testCases := []struct {
		name    string
		expPass bool
		req     *types.QueryNodeTotalRewardsRequest
		setup   func() *MockCaptains
	}{
		{
			name:    "fail - empty node id",
			expPass: false,
			req:     &types.QueryNodeTotalRewardsRequest{},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryNodeTotalRewards01)
			},
		},
		{
			name:    "fail - epoch == 1 epoch equal to 1",
			expPass: false,
			req:     &types.QueryNodeTotalRewardsRequest{NodeId: "node1"},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryNodeTotalRewards02)
			},
		},
		{
			name:    "pass - epoch == 2 and node with 0 rewards",
			expPass: true,
			req:     &types.QueryNodeTotalRewardsRequest{NodeId: "node1"},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryNodeTotalRewards03)
			},
		},
		{
			name:    "pass - epoch == 2 and node with 100 rewards",
			expPass: true,
			req:     &types.QueryNodeTotalRewardsRequest{NodeId: "node1"},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryNodeTotalRewards04)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("QueryNodeTotalRewards - %s", tc.name), func() {
			if tc.setup != nil {
				mockCaptains := tc.setup()
				suite.app.ClaimsKeeper.SetCaptainsKeeper(mockCaptains)
			}
			res, err := suite.queryClient.NodeTotalRewards(suite.ctx.Context(), tc.req)
			if tc.expPass {
				suite.Require().NotNil(res)
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})

	}
}

func (suite *ClaimsTestSuite) TestQueryHolderTotalRewards() {
	testCases := []struct {
		name    string
		expPass bool
		req     *types.QueryHolderTotalRewardsRequest
		setup   func() *MockCaptains
	}{
		{
			name:    "fail - empty owner address",
			expPass: false,
			req:     &types.QueryHolderTotalRewardsRequest{},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryHolderTotalRewards01)
			},
		},
		{
			name:    "fail - invalid owner address",
			expPass: false,
			req:     &types.QueryHolderTotalRewardsRequest{Owner: "bob"},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryHolderTotalRewards02)
			},
		},
		{
			name:    "fail - holder not found",
			expPass: false,
			req:     &types.QueryHolderTotalRewardsRequest{Owner: suite.cosmosAddress.String()},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryHolderTotalRewards03)
			},
		},
		{
			name:    "pass - epoch == 2 and owner has 1 node with 0 rewards",
			expPass: true,
			req:     &types.QueryHolderTotalRewardsRequest{Owner: suite.cosmosAddress.String()},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryHolderTotalRewards04)
			},
		},
		{
			name:    "pass - epoch == 2 and  owner has 1 node with 100 rewards",
			expPass: true,
			req:     &types.QueryHolderTotalRewardsRequest{Owner: suite.cosmosAddress.String()},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryHolderTotalRewards05)
			},
		},
		{
			name:    "pass - epoch == 2 and owner has 5 node with 0 rewards",
			expPass: true,
			req:     &types.QueryHolderTotalRewardsRequest{Owner: suite.cosmosAddress.String()},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryHolderTotalRewards06)
			},
		},
		{
			name:    "pass - epoch == 2 and owner has 5 node with 100 rewards",
			expPass: true,
			req:     &types.QueryHolderTotalRewardsRequest{Owner: suite.cosmosAddress.String()},
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyQueryHolderTotalRewards07)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("QueryHolderTotalRewards - %s", tc.name), func() {
			if tc.setup != nil {
				mockCaptains := tc.setup()
				suite.app.ClaimsKeeper.SetCaptainsKeeper(mockCaptains)
			}
			res, err := suite.queryClient.HolderTotalRewards(suite.ctx.Context(), tc.req)
			if tc.expPass {
				suite.Require().NotNil(res)
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})

	}

}
