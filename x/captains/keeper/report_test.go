package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *CaptainsTestSuite) TestHandleReportDigest() {

	testCases := []struct {
		name      string
		report    *types.ReportDigest
		expectErr bool
	}{
		{
			name: "success - handle report digest",
			report: &types.ReportDigest{
				EpochId:                1,
				GlobalOnOperationRatio: sdk.NewDecWithPrec(5, 1),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		err := suite.keeper.HandleReportDigest(suite.ctx, tc.report)
		if tc.expectErr {
			suite.Require().Errorf(err, "%s", err.Error())
		} else {
			suite.Require().NoError(err)
			val, _ := suite.keeper.GetEpochEmission(suite.ctx, tc.report.EpochId)
			suite.T().Log(val)
		}
	}
}

func (suite *CaptainsTestSuite) TestHandleReportBatch() {
	// prepare addresses and args
	addr1 := accounts[1].String()

	// prepare nodes
	nodes := suite.utilsBatchCreateCaptainNode(addr1, 1, 10)

	resp, err := suite.queryClient.Nodes(suite.ctx, &types.QueryNodesRequest{})
	suite.NoError(err)
	suite.Require().Len(resp.Nodes, len(nodes))

	ratios := suite.utilsBatchAssignFixedPowerOnRatio(nodes, 1, 0)

	// submit report digest on epoch 1
	err = suite.keeper.HandleReportDigest(suite.ctx, &types.ReportDigest{
		EpochId:                  suite.keeper.GetCurrentEpoch(suite.ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           3,
		MaximumNodeCountPerBatch: 4,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(5, 1),
	})
	suite.Require().NoError(err)

	// next height
	suite.Commit()
	suite.T().Logf("current height: %d", suite.ctx.BlockHeight())

	// submit report batch 1 on epoch 1
	err = suite.keeper.HandleReportBatch(suite.ctx, &types.ReportBatch{
		EpochId:   suite.keeper.GetCurrentEpoch(suite.ctx),
		BatchId:   1,
		NodeCount: 4,
		Nodes:     ratios[:4],
	})
	suite.Require().NoError(err)

	// next height
	suite.Commit()
	suite.T().Logf("current height: %d", suite.ctx.BlockHeight())

	// submit report batch 2 on epoch 1
	err = suite.keeper.HandleReportBatch(suite.ctx, &types.ReportBatch{
		EpochId:   suite.keeper.GetCurrentEpoch(suite.ctx),
		BatchId:   2,
		NodeCount: 4,
		Nodes:     ratios[4:8],
	})
	suite.Require().NoError(err)

	// next height
	suite.Commit()
	suite.T().Logf("current height: %d", suite.ctx.BlockHeight())

	// submit report batch 3 on epoch 1
	err = suite.keeper.HandleReportBatch(suite.ctx, &types.ReportBatch{
		EpochId:   suite.keeper.GetCurrentEpoch(suite.ctx),
		BatchId:   3,
		NodeCount: 2,
		Nodes:     ratios[8:],
	})
	suite.Require().NoError(err)

}
