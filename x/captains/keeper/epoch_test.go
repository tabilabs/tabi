package keeper_test

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

// In epoch(t), we check state value or calc transit value as follows:
//
// Before submitting digest:
//
// - pledge_sum(t)
// - historical_claimed_sum(t)
// - epoch_emission_base(t)
//
// After submitting digest:
//
// - epoch_emission_sum(t)
//
// Before submitting batch:
//
// - historical_emission_sum_by_node(t-2, n)
// - get: pledge_by_owner(t-1) (should be non-zero if has pledge)
//
// After submitting batch:
//
// - get: historical_emission_sum_by_node(t-1, n)
// - check: pledge_by_owner(t-1) (should be zero)
// - check: pledge_by_owner(t+1) (if has pledge)
//
//
// After submitting all batches:
//
// Before submitting end:
//
// After submitting end:
//
// Before enter new epoch:
//
// 1. PledgeSum
// 2.
// 3. ComputingPowerSum
// 4. NodeComputingPower
// 5.
//
// And we cal values as follows:
//

func (suite *CaptainsTestSuite) TestFullEpochPeriodWithoutOwnerStaking() {
	// create 100 nodes for addr1
	addr1 := accounts[1].String()
	nodes := suite.utilsBatchCreateCaptainNode(addr1, 1, 100)
	resp, err := suite.queryClient.Nodes(suite.ctx, &types.QueryNodesRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Nodes, len(nodes))

	///////////////////////////////////////////////
	//                 epoch 1                   //
	///////////////////////////////////////////////
	epoch1 := suite.keeper.GetCurrentEpoch(suite.ctx)
	suite.T().Logf("current epoch: %d", epoch1)

	// submit digest
	suite.T().Logf("subimit digest at height: %d", suite.ctx.BlockHeight())
	digest1 := types.ReportDigest{
		EpochId:                  suite.keeper.GetCurrentEpoch(suite.ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(1, 0),
	}
	anyVal, err := cdctypes.NewAnyWithValue(&digest1)
	suite.Require().NoError(err)

	_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)

	_, found := suite.keeper.GetDigest(suite.ctx, digest1.EpochId)
	suite.Require().True(found)
	suite.Commit()

	// submit batches
	nodeWithRatios1 := suite.utilsBatchAssignFixedPowerOnRatio(nodes, 1, 0)
	for i := uint64(1); i <= digest1.TotalBatchCount; i++ {
		batch := types.ReportBatch{
			EpochId:   epoch1,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios1[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		suite.T().Logf("subimit report-%d at height: %d", i, suite.ctx.BlockHeight())
		_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)

		suite.Commit()
	}

	// submit batch end
	end1 := types.ReportEnd{
		EpochId: suite.keeper.GetCurrentEpoch(suite.ctx),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end1)
	suite.Require().NoError(err)

	_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	suite.Commit()

	///////////////////////////////////////////////
	//                 epoch 2                   //
	///////////////////////////////////////////////

	// check epoch
	epoch2 := suite.keeper.GetCurrentEpoch(suite.ctx)
	suite.T().Logf("current epoch: %d", epoch2)
	suite.Require().Equal(epoch2, epoch1+1)

	// query results
	suite.T().Logf("computing power sum: %s", suite.keeper.GetComputingPowerSumOnEpoch(suite.ctx, epoch1).String())
	emission, _ := suite.keeper.GetEpochEmission(suite.ctx, epoch1)
	suite.T().Logf("epoch emission sum: %s", emission.String())

	// submit digest
	suite.T().Logf("subimit digest at height: %d", suite.ctx.BlockHeight())
	digest2 := types.ReportDigest{
		EpochId:                  suite.keeper.GetCurrentEpoch(suite.ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(47, 0),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&digest2)
	suite.Require().NoError(err)

	_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)

	_, found = suite.keeper.GetDigest(suite.ctx, digest2.EpochId)
	suite.Require().True(found)
	suite.Commit()

	// submit batches
	nodeWithRatios2 := suite.utilsBatchAssignRandomPowerOnRatio(nodes)
	for i := uint64(1); i <= digest2.TotalBatchCount; i++ {
		batch := types.ReportBatch{
			EpochId:   epoch2,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios2[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		suite.T().Logf("subimit report-%d at height: %d", i, suite.ctx.BlockHeight())
		_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)

		suite.Commit()
	}

	// submit batch end
	end2 := types.ReportEnd{
		EpochId: epoch2,
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end2)
	suite.Require().NoError(err)

	_, err = suite.msgServer.CommitReport(suite.ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	suite.Commit()

	///////////////////////////////////////////////
	//                 epoch 3                   //
	///////////////////////////////////////////////

	// check epoch
	epoch3 := suite.keeper.GetCurrentEpoch(suite.ctx)
	suite.T().Logf("current epoch: %d", epoch3)
	suite.Require().Equal(epoch3, epoch2+1)

	// query results
	suite.T().Logf("computing power sum: %s", suite.keeper.GetComputingPowerSumOnEpoch(suite.ctx, epoch2).String())
	emission, _ = suite.keeper.GetEpochEmission(suite.ctx, epoch2)
	suite.T().Logf("epoch emission sum: %s", emission.String())
}
