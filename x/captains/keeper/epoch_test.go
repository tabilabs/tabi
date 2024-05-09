package keeper_test

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *CaptainsTestSuite) TestFullEpoch() {
	// create 100 nodes for addr1
	addr1 := accounts[1].String()
	nodes := suite.utilsBatchCreateCaptainNode(addr1, 1, 100)
	resp, err := suite.queryClient.Nodes(suite.ctx, &types.QueryNodesRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Nodes, len(nodes))
	suite.T().Logf("current height: %d", suite.ctx.BlockHeight())

	// submit digest
	digest := types.ReportDigest{
		EpochId:                  suite.keeper.GetCurrentEpoch(suite.ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(1, 0),
	}
	anyVal, err := cdctypes.NewAnyWithValue(&digest)
	suite.Require().NoError(err)

	_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)

	_, found := suite.keeper.GetDigest(suite.ctx, digest.EpochId)
	suite.Require().True(found)
	suite.Commit()

	// submit batches
	nodeWithRatios := suite.utilsBatchAssignFixedPowerOnRatio(nodes, 1, 0)
	for i := uint64(1); i <= digest.TotalBatchCount; i++ {
		suite.T().Logf("current height: %d", suite.ctx.BlockHeight())

		batch := types.ReportBatch{
			EpochId:   suite.keeper.GetCurrentEpoch(suite.ctx),
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios[(i-1)*10 : i*10],
		}

		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})

		suite.Require().NoError(err)
	}

	// submit batch end

	// enter new epoch

	// check epoch

	// withdraw rewards

}
