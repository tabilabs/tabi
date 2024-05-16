package keeper_test

import "github.com/tabilabs/tabi/x/captains/types"

func (suite *CaptainsTestSuite) TestQueryNode() {
	testCases := []struct {
		name      string
		req       *types.QueryNodeRequest
		prepareFn func(*types.QueryNodeRequest) string
		expectErr bool
	}{
		{
			name: "success: node exist",
			req:  &types.QueryNodeRequest{},
			prepareFn: func(req *types.QueryNodeRequest) string {
				nodeId := suite.utilsCreateCaptainNode(
					accounts[1].String(),
					1,
				)
				req.NodeId = nodeId
				return nodeId
			},
			expectErr: false,
		},
		{
			name: "failure: node does not exist",
			req:  &types.QueryNodeRequest{},
			prepareFn: func(request *types.QueryNodeRequest) string {
				request.NodeId = "unknown-node-id"
				return ""
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			nodeId := tc.prepareFn(tc.req)

			resp, err := suite.queryClient.Node(suite.ctx, tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(nodeId, resp.Node.Id)
			}

			suite.SetupTest()
		})
	}
}

func (suite *CaptainsTestSuite) TestQueryNodes() {
	testCases := []struct {
		name      string
		req       *types.QueryNodesRequest
		prepareFn func() []string
		expectErr bool
	}{
		{
			name: "success: node exist",
			req:  &types.QueryNodesRequest{},
			prepareFn: func() []string {
				return suite.utilsBatchCreateCaptainNode(accounts[0].String(), 1, 10)
			},
			expectErr: false,
		},
		{
			name: "success: owner specified",
			req: &types.QueryNodesRequest{
				Owner: accounts[1].String(),
			},
			prepareFn: func() []string {
				suite.utilsBatchCreateCaptainNode(accounts[0].String(), 1, 10)
				return suite.utilsBatchCreateCaptainNode(accounts[1].String(), 1, 5)
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			nodeIds := tc.prepareFn()

			resp, err := suite.queryClient.Nodes(suite.ctx, tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Len(resp.Nodes, len(nodeIds))
			}

			suite.SetupTest()
		})

	}
}

func (suite *CaptainsTestSuite) TestNodeLastEpochInfo() {
	// TODO: add tests
}

func (suite *CaptainsTestSuite) TestQueryDivision() {
	testCases := []struct {
		name        string
		req         *types.QueryDivisionRequest
		expectLevel uint64
		expectErr   bool
	}{
		{
			name: "success: division 1",
			req: &types.QueryDivisionRequest{
				DivisionId: types.GenDivisionsId(1),
			},
			expectLevel: 1,
			expectErr:   false,
		},
		{
			name: "success: division 2",
			req: &types.QueryDivisionRequest{
				DivisionId: types.GenDivisionsId(2),
			},
			expectLevel: 2,
			expectErr:   false,
		},
		{
			name: "success: division 3",
			req: &types.QueryDivisionRequest{
				DivisionId: types.GenDivisionsId(3),
			},
			expectLevel: 3,
			expectErr:   false,
		},
		{
			name: "success: division 4",
			req: &types.QueryDivisionRequest{
				DivisionId: types.GenDivisionsId(4),
			},
			expectLevel: 4,
			expectErr:   false,
		},
		{
			name: "success: division 5",
			req: &types.QueryDivisionRequest{
				DivisionId: types.GenDivisionsId(5),
			},
			expectLevel: 5,
			expectErr:   false,
		},
		{
			name: "failure: division does not exist",
			req: &types.QueryDivisionRequest{
				DivisionId: "unknown-division-id",
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.queryClient.Division(suite.ctx, tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expectLevel, resp.Division.Level)
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestQueryDivisions() {
	testCases := []struct {
		name      string
		req       *types.QueryDivisionsRequest
		expectErr bool
	}{
		{
			name:      "success: all divisions",
			req:       &types.QueryDivisionsRequest{},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.queryClient.Divisions(suite.ctx, tc.req)
			suite.Require().NoError(err)
			suite.Require().Len(resp.Divisions, 5)
		})
	}
}

func (suite *CaptainsTestSuite) TestQuerySupply() {
	testCases := []struct {
		name         string
		melleateFn   func()
		req          *types.QuerySupplyRequest
		expectSupply uint64
		expectErr    bool
	}{
		{
			name: "success: division 1",
			melleateFn: func() {
				suite.utilsBatchCreateCaptainNode(accounts[0].String(), 1, 10)
			},
			req: &types.QuerySupplyRequest{
				DivisionId: types.GenDivisionsId(1),
			},
			expectSupply: 10,
			expectErr:    false,
		},
		{
			name: "success: division 2",
			req: &types.QuerySupplyRequest{
				DivisionId: types.GenDivisionsId(2),
			},
			melleateFn: func() {
				suite.utilsBatchCreateCaptainNode(accounts[0].String(), 2, 10)
			},
			expectSupply: 10,
			expectErr:    false,
		},
		{
			name: "success: division 3",
			req: &types.QuerySupplyRequest{
				DivisionId: types.GenDivisionsId(3),
			},
			melleateFn: func() {
				suite.utilsBatchCreateCaptainNode(accounts[0].String(), 3, 10)
			},
			expectSupply: 10,
			expectErr:    false,
		},
		{
			name: "success: division 4",
			req: &types.QuerySupplyRequest{
				DivisionId: types.GenDivisionsId(4),
			},
			melleateFn: func() {
				suite.utilsBatchCreateCaptainNode(accounts[0].String(), 4, 10)
			},
			expectSupply: 10,
			expectErr:    false,
		},
		{
			name: "success: division 5",
			req: &types.QuerySupplyRequest{
				DivisionId: types.GenDivisionsId(5),
			},
			melleateFn: func() {
				suite.utilsBatchCreateCaptainNode(accounts[0].String(), 5, 10)
			},
			expectSupply: 10,
			expectErr:    false,
		},
		{
			name: "failure: unknown division",
			req: &types.QuerySupplyRequest{
				DivisionId: types.GenDivisionsId(6),
			},
			melleateFn: func() {},
			expectErr:  true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.melleateFn()

			resp, err := suite.queryClient.Supply(suite.ctx, tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expectSupply, resp.Amount)
			}

			suite.SetupTest()
		})
	}
}

func (suite *CaptainsTestSuite) TestQuerySaleLevel() {
	testCases := []struct {
		name        string
		req         *types.QuerySaleLevelRequest
		expectLevel uint64
		expectErr   bool
	}{
		{
			name:        "success: sale level 1",
			req:         &types.QuerySaleLevelRequest{},
			expectLevel: 1,
			expectErr:   false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.queryClient.SaleLevel(suite.ctx, tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().Equal(tc.expectLevel, resp.SaleLevel)
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestAuthorizedMembers() {
	testCases := []struct {
		name      string
		req       *types.QueryAuthorizedMembersRequest
		expectLen uint64
		expectErr bool
	}{
		{
			name:      "success",
			req:       &types.QueryAuthorizedMembersRequest{},
			expectLen: 1,
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.queryClient.AuthorizedMembers(suite.ctx, tc.req)

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Len(resp.Members, int(tc.expectLen))
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestCurrentEpoch() {
	testCases := []struct {
		name        string
		req         *types.QueryCurrentEpochRequest
		expectEpoch uint64
		expectErr   bool
	}{
		{
			name:        "success",
			req:         &types.QueryCurrentEpochRequest{},
			expectEpoch: 1,
			expectErr:   false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.queryClient.CurrentEpoch(suite.ctx, tc.req)

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expectEpoch, resp.Epoch)
			}
		})
	}
}
