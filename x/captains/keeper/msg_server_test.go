package keeper_test

import (
	"fmt"

	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *CaptainsTestSuite) TestCreateCaptainNode() {
	member := accounts[0].String()
	owner := accounts[1].String()

	// query divisions
	divisions := suite.utilsGetDivisions()

	testCases := []struct {
		name      string
		request   *types.MsgCreateCaptainNode
		expectErr bool
	}{
		{
			name: "success - create with authority",
			request: &types.MsgCreateCaptainNode{
				Authority:  member,
				Owner:      owner,
				DivisionId: divisions[1],
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgCreateCaptainNode - %s", tc.name), func() {
			_, err := suite.msgServer.CreateCaptainNode(suite.ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestCommitReport() {
	panic("implement me")
}

func (suite *CaptainsTestSuite) TestAddAuthorizedMembers() {
	member := accounts[0].String()
	owner := accounts[1].String()

	testCases := []struct {
		name      string
		request   *types.MsgAddAuthorizedMembers
		expectErr bool
	}{
		{
			name: "success - create with authority",
			request: &types.MsgAddAuthorizedMembers{
				Authority: member,
				Members:   []string{owner},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgAddAuthorizedMembers - %s", tc.name), func() {
			_, err := suite.msgServer.AddAuthorizedMembers(suite.ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestRemoveAuthorizedMembers() {
	member := accounts[0].String()
	toBeMember := accounts[1].String()
	nonMember := accounts[2].String()

	testCases := []struct {
		name      string
		malleate  func()
		request   *types.MsgRemoveAuthorizedMembers
		expectErr bool
	}{
		{
			name: "success - remove with authority",
			malleate: func() {
				suite.utilsAddAuthorizedMember(toBeMember)
			},
			request: &types.MsgRemoveAuthorizedMembers{
				Authority: member,
				Members:   []string{toBeMember},
			},
			expectErr: false,
		},
		{
			name: "failure - remove without authority",
			malleate: func() {
				suite.utilsAddAuthorizedMember(toBeMember)
			},
			request: &types.MsgRemoveAuthorizedMembers{
				Authority: nonMember,
				Members:   []string{toBeMember},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgRemoveAuthorizedMembers - %s", tc.name), func() {
			tc.malleate()
			_, err := suite.msgServer.RemoveAuthorizedMembers(suite.ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestUpdateSaleLevel() {
	member := accounts[0].String()
	nonMember := accounts[1].String()

	testCases := []struct {
		name      string
		request   *types.MsgUpdateSaleLevel
		expectErr bool
	}{
		{
			name: "success - update with authority",
			request: &types.MsgUpdateSaleLevel{
				Authority: member,
				SaleLevel: 2,
			},
			expectErr: false,
		},
		{
			name: "failure - update without authority",
			request: &types.MsgUpdateSaleLevel{
				Authority: nonMember,
				SaleLevel: 2,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgUpdateSaleLevel - %s", tc.name), func() {
			_, err := suite.msgServer.UpdateSaleLevel(suite.ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestCommitComputingPower() {
	member := accounts[0].String()
	owner := accounts[1].String()

	// prepare a node
	suite.utilsCreateCaptainNode(owner, 1)

	testCases := []struct {
		name      string
		request   *types.MsgCommitComputingPower
		expectErr bool
	}{
		{
			name: "success - commit with authority",
			request: &types.MsgCommitComputingPower{
				Authority: member,
				ComputingPowerRewards: []types.ClaimableComputingPower{
					{10000, owner},
				},
			},
			expectErr: false,
		},
		{
			name: "failure - commit without authority",
			request: &types.MsgCommitComputingPower{
				Authority: owner,
				ComputingPowerRewards: []types.ClaimableComputingPower{
					{10000, owner},
				},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgCommitComputingPower - %s", tc.name), func() {
			_, err := suite.msgServer.CommitComputingPower(suite.ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *CaptainsTestSuite) TestClaimComputingPower() {
	// prepare addresses
	owner := accounts[1].String()

	// create nodes & commit powers
	nodeID := suite.utilsCreateCaptainNode(owner, 1)
	suite.utilsCommitPower(owner, 10000)

	_, err := suite.queryClient.Node(suite.ctx, &types.QueryNodeRequest{NodeId: nodeID})
	suite.Require().NoError(err)

	testCases := []struct {
		name      string
		malleate  func()
		request   *types.MsgClaimComputingPower
		expectErr bool
	}{
		{
			name: "success - claim power enough",
			malleate: func() {
			},
			request: &types.MsgClaimComputingPower{
				Sender:               owner,
				ComputingPowerAmount: 10000,
				NodeId:               nodeID,
			},
			expectErr: false,
		},
		{
			name: "failure - claim power exceeded",
			malleate: func() {
			},
			request: &types.MsgClaimComputingPower{
				Sender:               owner,
				ComputingPowerAmount: 100000,
				NodeId:               nodeID,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgClaimComputingPower - %s", tc.name), func() {
			tc.malleate()
			_, err := suite.msgServer.ClaimComputingPower(suite.ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
