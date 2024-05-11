package keeper_test

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

//  In epoch(t), we check state value or calc intermediate value as follows:
//
//	A1. BeforeReportDigest(T)
//
//	| Name                    | Action | Expect                                  |
//	| ----------------------- | ------ | --------------------------------------- |
//	| PledgeSum(T)            | Get    | - Staked: NonZero<br>- Unstaked: Zero   |
//	| HistoricalClaimedSum    | Get    | - Claimed: NonZero<br>- Unclaimed: Zero |
//	| EpochEmissionSum(T)@Ref | Calc   | For comparison                          |
//
//	A2. AfterReportDigest(T)
//
//	| Name                | Action | Expect                  |
//	| ------------------- | ------ | ----------------------- |
//	| PledgeSum(T)        | Get    | Zero                    |
//	| EpochEmissionSum(T) | Get    | EpochEmissionSum(T)@Ref |
//
//	B1. BeforeReportBatch(T)
//
//	| Name                              | Action | Expect                         |
//	| --------------------------------- | ------ | ------------------------------ |
//	| HistoricalEmissionByNode(T-2)     | Get    | - T<=2, Zero<br>- T>2, NonZero |
//	| EpochEmissionSum(T-1)             | Get    | - T=1, Zero<br>- T=2, NonZero  |
//	| ComputingPowerSum(T-1)            | Get    |                                |
//	| ComputingPowerByNode(T-1)         | Get    |                                |
//	| HistoricalEmissionByNode(T-1)@Ref | Calc   | For comparison                 |
//	| PledgeByOwner(T-2)                | Get    |                                |
//	| ComputingPowerByNode(T)@Ref       | Calc   | For comparison                 |
//	| ComputingPowerSum(T)@Ref          | Calc   | For comparison                 |
//	| PledgeByOwner(T+1)                | Get    | Zero                           |
//
//	B2. AfterReportBatch(T)
//
//	| Name                          | Action | Expect                            |
//	| ----------------------------- | ------ | --------------------------------- |
//	| HistoricalEmissionByNode(T-2) | Get    | deleted                           |
//	| ComputingPower(T-2)           | Get    | deleted                           |
//	| HistoricalEmissionByNode(T-1) | Get    | HistoricalEmissionByNode(T-1)@Ref |
//	| PledgeByOwner(T-2)            | Get    | deleted                           |
//	| ComputingPowerByNode(T)       | Get    | ComputingPowerByNode(T)@Ref       |
//	| ComputingPowerSum(T)          | Get    | ComputingPowerSum(T)@Ref          |
//	| PledgeByOwner(T+1)            | Get    | NonZero if staked                 |
//
//	C1. AfterEpochEnd(T)
//
//	| Name                   | Action | Expect   |
//	| ---------------------- | ------ | -------- |
//	| EpochEmissionSum(T-1)  | Get    | Non zero |
//	| ComputingPowerSum(T-1) | Get    | Non zero |
//
//	D1. BeginBlocker
//
//	| Name                   | Action | Expect  |
//	| ---------------------- | ------ | ------- |
//	| EpochEmissionSum(T-1)  | Get    | deleted |
//	| ComputingPowerSum(T-1) | Get    | deleted |

// TestFullEpochPeriodWithoutOwnerStaking tests the full epoch period without owner staking or claiming.
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

	// BeforeReportDigest(1)
	//

	// Submit Report Digest
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

	// AfterReportDigest(1)
	emission := suite.keeper.GetEpochEmission(suite.ctx, epoch1)
	suite.T().Logf("emission at 1: %s", emission.String())

	// Before Submit Report Batch

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
	emission = suite.keeper.GetEpochEmission(suite.ctx, epoch2)
	suite.T().Logf("epoch emission sum: %s", emission.String())
}
