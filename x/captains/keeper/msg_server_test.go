package keeper_test

import (
	"fmt"

	sdkcdc "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *IntegrationTestSuite) TestCreateCaptainNode() {
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
			_, err := suite.MsgServer.CreateCaptainNode(suite.Ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestAddAuthorizedMembers() {
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
			_, err := suite.MsgServer.AddAuthorizedMembers(suite.Ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestRemoveAuthorizedMembers() {
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
			_, err := suite.MsgServer.RemoveAuthorizedMembers(suite.Ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestUpdateSaleLevel() {
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
			_, err := suite.MsgServer.UpdateSaleLevel(suite.Ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestCommitComputingPower() {
	// NOTE: only test report validation here, the full commit logic is tested
	// in epoch_test.go
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
			_, err := suite.MsgServer.CommitComputingPower(suite.Ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestClaimComputingPower() {
	// prepare addresses
	authority := accounts[0].String()
	owner := accounts[1].String()

	// create nodes & commit powers
	nodeID := suite.utilsCreateCaptainNode(owner, 1)
	suite.utilsCommitPower(owner, 10000)

	_, err := suite.QueryClient.Node(suite.Ctx, &types.QueryNodeRequest{NodeId: nodeID})
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
				Sender:               authority,
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
				Sender:               authority,
				ComputingPowerAmount: 100000,
				NodeId:               nodeID,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgClaimComputingPower - %s", tc.name), func() {
			tc.malleate()
			_, err := suite.MsgServer.ClaimComputingPower(suite.Ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestCommitReport() {
	member := accounts[0].String()
	owner := accounts[1].String()

	testCases := []struct {
		name          string
		mellateReport func(*types.MsgCommitReport)
		request       *types.MsgCommitReport
		expectErr     bool
	}{
		{
			name:          "failure - unauthorized member",
			mellateReport: func(msg *types.MsgCommitReport) {},
			request: &types.MsgCommitReport{
				Authority: owner,
			},
			expectErr: true,
		},
		{
			name:          "failure - invalid report type",
			mellateReport: func(msg *types.MsgCommitReport) {},
			request: &types.MsgCommitReport{
				Authority:  member,
				ReportType: types.ReportType_REPORT_TYPE_UNSPECIFIED,
			},
			expectErr: true,
		},
		{
			name: "failure - invalid report digest epoch",
			mellateReport: func(msg *types.MsgCommitReport) {
				report := types.ReportDigest{
					EpochId:                  0,
					TotalBatchCount:          100,
					TotalNodeCount:           10000,
					MaximumNodeCountPerBatch: 10,
					GlobalOnOperationRatio:   sdk.Dec{},
				}
				val, err := sdkcdc.NewAnyWithValue(&report)
				if err != nil {
					panic(err)
				}
				msg.Report = val
			},
			request: &types.MsgCommitReport{
				Authority:  member,
				ReportType: types.ReportType_REPORT_TYPE_DIGEST,
			},
			expectErr: true,
		},
		{
			name: "failure - digest not found on epoch",
			mellateReport: func(msg *types.MsgCommitReport) {
				report := types.ReportBatch{
					EpochId: 1,
					BatchId: 10,
				}
				val, err := sdkcdc.NewAnyWithValue(&report)
				if err != nil {
					panic(err)
				}
				msg.Report = val
			},
			request: &types.MsgCommitReport{
				Authority:  member,
				ReportType: types.ReportType_REPORT_TYPE_BATCH,
			},
			expectErr: true,
		},
		{
			name: "failure - invalid report end fields",
			mellateReport: func(msg *types.MsgCommitReport) {
				report := types.ReportEnd{
					EpochId: 100,
				}
				val, err := sdkcdc.NewAnyWithValue(&report)
				if err != nil {
					panic(err)
				}
				msg.Report = val
			},
			request: &types.MsgCommitReport{
				Authority:  member,
				ReportType: types.ReportType_REPORT_TYPE_END,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgCommitReport - %s", tc.name), func() {
			tc.mellateReport(tc.request)

			_, err := suite.MsgServer.CommitReport(suite.Ctx, tc.request)

			if tc.expectErr {
				suite.Require().Errorf(err, "%s", err.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
