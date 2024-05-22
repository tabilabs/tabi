package keeper_test

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	claimstypes "github.com/tabilabs/tabi/x/claims/types"

	"github.com/tabilabs/tabi/x/captains/types"
)

//  In epoch(t), check state value or calc intermediate value as follows:
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
//	D1. BeginBlocker(T+1)
//
//	| Name                   | Action | Expect  |
//	| ---------------------- | ------ | ------- |
//	| EpochEmissionSum(T-1)  | Get    | deleted |
//	| ComputingPowerSum(T-1) | Get    | deleted |

// TestCompletedEpochesScenario1 tests the full epoch period but:
// 1. no user claim their rewards: so pledge ratio is always 1
// 2. no user stake their coins: so not historical claimed emissions
func (suite *IntegrationTestSuite) TestCompletedEpochesScenario1() {
	// create 100 nodes for addr1
	addr1 := accounts[1].String()
	nodes := suite.utilsBatchCreateCaptainNode(addr1, 1, 100)
	resp, err := suite.QueryClient.Nodes(suite.Ctx, &types.QueryNodesRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Nodes, len(nodes))
	nodeWithRatios := suite.utilsBatchAssignFixedPowerOnRatio(nodes, 1, 0)

	//////////////////////////////////////////////////////////////////////
	//                                epoch 1                           //
	//////////////////////////////////////////////////////////////////////
	epoch1 := suite.Keeper.GetCurrentEpoch(suite.Ctx)
	// BeforeReportDigest(1)
	expectEmission1 := suite.Keeper.CalcEpochEmission(suite.Ctx, epoch1, sdk.NewDecWithPrec(1, 0))
	// Submit ReportDigest(1)
	digest1 := types.ReportDigest{
		EpochId:                  suite.Keeper.GetCurrentEpoch(suite.Ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(1, 0),
	}
	anyVal, err := cdctypes.NewAnyWithValue(&digest1)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)
	suite.Commit()

	// AfterReportDigest(1)
	_, found := suite.Keeper.GetReportDigest(suite.Ctx, digest1.EpochId)
	suite.Require().True(found)
	actualEmission1 := suite.Keeper.GetEpochEmission(suite.Ctx, epoch1)
	suite.Require().Equal(expectEmission1, actualEmission1)

	// Submit ReportBatches(1)
	expectedComputingPowerSum1 := sdk.NewDec(0)
	for i := uint64(1); i <= digest1.TotalBatchCount; i++ {
		// Before Submit ReportBatches(1,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			expectedNodeComputingPower := suite.Keeper.CalcNodeComputingPowerOnEpoch(suite.Ctx,
				epoch1, node.NodeId, node.OnOperationRatio)
			expectedComputingPowerSum1 = expectedComputingPowerSum1.Add(expectedNodeComputingPower)
		}
		// Submit ReportBatch(1, i)
		batch := types.ReportBatch{
			EpochId:   epoch1,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)
		suite.Commit()
		// After Submit ReportBatches(1,i)
		actualComputingPowerSum1 := suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1)
		suite.Require().Equal(expectedComputingPowerSum1, actualComputingPowerSum1)
	}

	// Submit ReportEnd(1)
	end1 := types.ReportEnd{
		EpochId: suite.Keeper.GetCurrentEpoch(suite.Ctx),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end1)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	suite.Commit()

	//////////////////////////////////////////////////////////////////////
	//                                epoch 2                           //
	//////////////////////////////////////////////////////////////////////
	epoch2 := suite.Keeper.GetCurrentEpoch(suite.Ctx)
	// BeforeReportDigest(2)
	expectEmission2 := suite.Keeper.CalcEpochEmission(suite.Ctx, epoch2, sdk.NewDecWithPrec(47, 1))
	// Submit ReportDigest(2)
	digest2 := types.ReportDigest{
		EpochId:                  suite.Keeper.GetCurrentEpoch(suite.Ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(47, 1),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&digest2)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)
	suite.Commit()

	// AfterReportDigest(2)
	_, found = suite.Keeper.GetReportDigest(suite.Ctx, digest2.EpochId)
	suite.Require().True(found)
	actualEmission2 := suite.Keeper.GetEpochEmission(suite.Ctx, epoch2)
	suite.Require().Equal(expectEmission2, actualEmission2)

	// Submit ReportBatches(2)
	suite.Require().Equal(expectEmission1, suite.Keeper.GetEpochEmission(suite.Ctx, epoch1))
	suite.Require().Equal(expectedComputingPowerSum1, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1))

	expectedComputingPowerSum2 := sdk.NewDec(0)
	for i := uint64(1); i <= digest2.TotalBatchCount; i++ {
		// Before Submit ReportBatches(2,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			expectedNodeComputingPower := suite.Keeper.CalcNodeComputingPowerOnEpoch(suite.Ctx,
				epoch2, node.NodeId, node.OnOperationRatio)
			expectedComputingPowerSum2 = expectedComputingPowerSum2.Add(expectedNodeComputingPower)
		}
		// Submit ReportBatch(2, i)
		batch := types.ReportBatch{
			EpochId:   epoch2,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)
		suite.Commit()
		// After Submit ReportBatches(2,i)
		actualComputingPowerSum2 := suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2)
		suite.Require().Equal(expectedComputingPowerSum2, actualComputingPowerSum2)
	}

	// Submit ReportEnd(2)
	end2 := types.ReportEnd{
		EpochId: suite.Keeper.GetCurrentEpoch(suite.Ctx),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end2)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	// After ReportEnd(2)
	suite.Require().Equal(actualEmission1, suite.Keeper.GetEpochEmission(suite.Ctx, epoch1))
	suite.Require().Equal(expectedComputingPowerSum1, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1))

	suite.Commit()
	// BeginBlocker
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetEpochEmission(suite.Ctx, epoch1))
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1))

	//////////////////////////////////////////////////////////////////////
	//                                epoch 3                           //
	//////////////////////////////////////////////////////////////////////
	epoch3 := suite.Keeper.GetCurrentEpoch(suite.Ctx)
	// BeforeReportDigest(3)
	expectEmission3 := suite.Keeper.CalcEpochEmission(suite.Ctx, epoch3, sdk.NewDecWithPrec(66, 1))
	// Submit ReportDigest(2)
	digest3 := types.ReportDigest{
		EpochId:                  epoch3,
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(66, 1),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&digest3)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)
	suite.Commit()

	// AfterReportDigest(3)
	_, found = suite.Keeper.GetReportDigest(suite.Ctx, digest3.EpochId)
	suite.Require().True(found)
	actualEmission3 := suite.Keeper.GetEpochEmission(suite.Ctx, epoch3)
	suite.Require().Equal(expectEmission3, actualEmission3)

	// Submit ReportBatches(3)
	suite.Require().Equal(expectEmission2, suite.Keeper.GetEpochEmission(suite.Ctx, epoch2))
	suite.Require().Equal(expectedComputingPowerSum2, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2))

	expectedComputingPowerSum3 := sdk.NewDec(0)
	for i := uint64(1); i <= digest2.TotalBatchCount; i++ {
		// Before Submit ReportBatches(3,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			suite.Require().NotEmpty(suite.Keeper.GetNodeCumulativeEmissionByEpoch(suite.Ctx, epoch1, node.NodeId))
			suite.Require().NotEmpty(suite.Keeper.GetNodeComputingPowerOnEpoch(suite.Ctx, epoch2, node.NodeId))

			expectedNodeComputingPower := suite.Keeper.CalcNodeComputingPowerOnEpoch(suite.Ctx,
				epoch3, node.NodeId, node.OnOperationRatio)
			expectedComputingPowerSum3 = expectedComputingPowerSum3.Add(expectedNodeComputingPower)
		}
		// Submit ReportBatch(3, i)
		batch := types.ReportBatch{
			EpochId:   epoch3,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)
		suite.Commit()

		// After Submit ReportBatches(2,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetNodeCumulativeEmissionByEpoch(suite.Ctx, epoch1, node.NodeId))
			suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetNodeComputingPowerOnEpoch(suite.Ctx, epoch1, node.NodeId))
		}

		actualComputingPowerSum3 := suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch3)
		suite.Require().Equal(expectedComputingPowerSum3, actualComputingPowerSum3)
	}

	// Submit ReportEnd(3)
	end3 := types.ReportEnd{
		EpochId: suite.Keeper.GetCurrentEpoch(suite.Ctx),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end3)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	// After ReportEnd(3)
	suite.Require().Equal(actualEmission2, suite.Keeper.GetEpochEmission(suite.Ctx, epoch2))
	suite.Require().Equal(expectedComputingPowerSum2, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2))

	suite.Commit()
	// BeginBlocker
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetEpochEmission(suite.Ctx, epoch2))
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2))
}

// TestCompletedEpochesScenario1 tests the full epoch period but:
// 1. user claim their rewards: so pledge ratio is calculated
func (suite *IntegrationTestSuite) TestCompletedEpochesScenario2() {
	// create 100 nodes for addr1
	addr1 := accounts[1].String()
	nodes := suite.utilsBatchCreateCaptainNode(addr1, 1, 100)
	resp, err := suite.QueryClient.Nodes(suite.Ctx, &types.QueryNodesRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Nodes, len(nodes))
	nodeWithRatios := suite.utilsBatchAssignFixedPowerOnRatio(nodes, 1, 0)

	//////////////////////////////////////////////////////////////////////
	//                                epoch 1                           //
	//////////////////////////////////////////////////////////////////////
	epoch1 := suite.Keeper.GetCurrentEpoch(suite.Ctx)
	// BeforeReportDigest(1)
	expectEmission1 := suite.Keeper.CalcEpochEmission(suite.Ctx, epoch1, sdk.NewDecWithPrec(1, 0))
	// Submit ReportDigest(1)
	digest1 := types.ReportDigest{
		EpochId:                  suite.Keeper.GetCurrentEpoch(suite.Ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(1, 0),
	}
	anyVal, err := cdctypes.NewAnyWithValue(&digest1)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)
	suite.Commit()

	// AfterReportDigest(1)
	_, found := suite.Keeper.GetReportDigest(suite.Ctx, digest1.EpochId)
	suite.Require().True(found)
	actualEmission1 := suite.Keeper.GetEpochEmission(suite.Ctx, epoch1)
	suite.Require().Equal(expectEmission1, actualEmission1)
	suite.T().Logf("emission at %d: %s", epoch1, actualEmission1)

	// Submit ReportBatches(1)
	expectedComputingPowerSum1 := sdk.NewDec(0)
	for i := uint64(1); i <= digest1.TotalBatchCount; i++ {
		// Before Submit ReportBatches(1,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			expectedNodeComputingPower := suite.Keeper.CalcNodeComputingPowerOnEpoch(suite.Ctx,
				epoch1, node.NodeId, node.OnOperationRatio)
			expectedComputingPowerSum1 = expectedComputingPowerSum1.Add(expectedNodeComputingPower)
		}
		// Submit ReportBatch(1, i)
		batch := types.ReportBatch{
			EpochId:   epoch1,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)
		suite.Commit()
		// After Submit ReportBatches(1,i)
		actualComputingPowerSum1 := suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1)
		suite.Require().Equal(expectedComputingPowerSum1, actualComputingPowerSum1)
	}
	suite.T().Logf("computing power sum at %d: %s", epoch1, expectedComputingPowerSum1)

	// Submit ReportEnd(1)
	end1 := types.ReportEnd{
		EpochId: suite.Keeper.GetCurrentEpoch(suite.Ctx),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end1)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	suite.Commit()

	//////////////////////////////////////////////////////////////////////
	//                                epoch 2                           //
	//////////////////////////////////////////////////////////////////////
	epoch2 := suite.Keeper.GetCurrentEpoch(suite.Ctx)

	// NOTE: claim rewards
	resp1, err := suite.ClaimsServer.Claims(suite.Ctx, &claimstypes.MsgClaims{
		Receiver: addr1,
		Sender:   addr1,
	})
	suite.Require().NoError(err)
	suite.T().Logf("claimed rewards at %d: %s", epoch2, resp1.Amount)

	// BeforeReportDigest(2)
	expectEmission2 := suite.Keeper.CalcEpochEmission(suite.Ctx, epoch2, sdk.NewDecWithPrec(47, 1))
	// Submit ReportDigest(2)
	digest2 := types.ReportDigest{
		EpochId:                  suite.Keeper.GetCurrentEpoch(suite.Ctx),
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(47, 1),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&digest2)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)
	suite.Commit()

	// AfterReportDigest(2)
	_, found = suite.Keeper.GetReportDigest(suite.Ctx, digest2.EpochId)
	suite.Require().True(found)
	actualEmission2 := suite.Keeper.GetEpochEmission(suite.Ctx, epoch2)
	suite.Require().Equal(expectEmission2, actualEmission2)
	suite.T().Logf("emission at %d: %s", epoch2, actualEmission2)

	// Submit ReportBatches(2)
	suite.Require().Equal(expectEmission1, suite.Keeper.GetEpochEmission(suite.Ctx, epoch1))
	suite.Require().Equal(expectedComputingPowerSum1, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1))

	expectedComputingPowerSum2 := sdk.NewDec(0)
	for i := uint64(1); i <= digest2.TotalBatchCount; i++ {
		// Before Submit ReportBatches(2,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			expectedNodeComputingPower := suite.Keeper.CalcNodeComputingPowerOnEpoch(suite.Ctx,
				epoch2, node.NodeId, node.OnOperationRatio)
			expectedComputingPowerSum2 = expectedComputingPowerSum2.Add(expectedNodeComputingPower)
		}
		// Submit ReportBatch(2, i)
		batch := types.ReportBatch{
			EpochId:   epoch2,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)
		suite.Commit()
		// After Submit ReportBatches(2,i)
		actualComputingPowerSum2 := suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2)
		suite.Require().Equal(expectedComputingPowerSum2, actualComputingPowerSum2)
	}

	// Submit ReportEnd(2)
	end2 := types.ReportEnd{
		EpochId: suite.Keeper.GetCurrentEpoch(suite.Ctx),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end2)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	// After ReportEnd(2)
	suite.Require().Equal(actualEmission1, suite.Keeper.GetEpochEmission(suite.Ctx, epoch1))
	suite.Require().Equal(expectedComputingPowerSum1, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1))

	suite.Commit()
	// BeginBlocker
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetEpochEmission(suite.Ctx, epoch1))
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch1))

	//////////////////////////////////////////////////////////////////////
	//                                epoch 3                           //
	//////////////////////////////////////////////////////////////////////
	epoch3 := suite.Keeper.GetCurrentEpoch(suite.Ctx)
	// BeforeReportDigest(3)
	expectEmission3 := suite.Keeper.CalcEpochEmission(suite.Ctx, epoch3, sdk.NewDecWithPrec(66, 1))
	// Submit ReportDigest(2)
	digest3 := types.ReportDigest{
		EpochId:                  epoch3,
		TotalBatchCount:          10,
		TotalNodeCount:           100,
		MaximumNodeCountPerBatch: 10,
		GlobalOnOperationRatio:   sdk.NewDecWithPrec(66, 1),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&digest3)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		Report:     anyVal,
		ReportType: types.ReportType_REPORT_TYPE_DIGEST,
	})
	suite.Require().NoError(err)
	suite.Commit()

	// AfterReportDigest(3)
	_, found = suite.Keeper.GetReportDigest(suite.Ctx, digest3.EpochId)
	suite.Require().True(found)
	actualEmission3 := suite.Keeper.GetEpochEmission(suite.Ctx, epoch3)
	suite.Require().Equal(expectEmission3, actualEmission3)
	suite.T().Logf("emission at %d: %s", epoch3, actualEmission3)

	// Submit ReportBatches(3)
	suite.Require().Equal(expectEmission2, suite.Keeper.GetEpochEmission(suite.Ctx, epoch2))
	suite.Require().Equal(expectedComputingPowerSum2, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2))

	expectedComputingPowerSum3 := sdk.NewDec(0)
	for i := uint64(1); i <= digest2.TotalBatchCount; i++ {
		// Before Submit ReportBatches(3,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			suite.Require().NotEmpty(suite.Keeper.GetNodeCumulativeEmissionByEpoch(suite.Ctx, epoch1, node.NodeId))
			suite.Require().NotEmpty(suite.Keeper.GetNodeComputingPowerOnEpoch(suite.Ctx, epoch2, node.NodeId))

			expectedNodeComputingPower := suite.Keeper.CalcNodeComputingPowerOnEpoch(suite.Ctx,
				epoch3, node.NodeId, node.OnOperationRatio)
			expectedComputingPowerSum3 = expectedComputingPowerSum3.Add(expectedNodeComputingPower)
		}
		// Submit ReportBatch(3, i)
		batch := types.ReportBatch{
			EpochId:   epoch3,
			BatchId:   i,
			NodeCount: 10,
			Nodes:     nodeWithRatios[(i-1)*10 : i*10],
		}
		anyVal, err := cdctypes.NewAnyWithValue(&batch)
		suite.Require().NoError(err)

		_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
			Authority:  accounts[0].String(),
			Report:     anyVal,
			ReportType: types.ReportType_REPORT_TYPE_BATCH,
		})
		suite.Require().NoError(err)
		suite.Commit()

		// After Submit ReportBatches(2,i)
		for _, node := range nodeWithRatios[(i-1)*10 : i*10] {
			suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetNodeCumulativeEmissionByEpoch(suite.Ctx, epoch1, node.NodeId))
			suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetNodeComputingPowerOnEpoch(suite.Ctx, epoch1, node.NodeId))
		}

		actualComputingPowerSum3 := suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch3)
		suite.Require().Equal(expectedComputingPowerSum3, actualComputingPowerSum3)
	}

	// Submit ReportEnd(3)
	end3 := types.ReportEnd{
		EpochId: suite.Keeper.GetCurrentEpoch(suite.Ctx),
	}
	anyVal, err = cdctypes.NewAnyWithValue(&end3)
	suite.Require().NoError(err)

	_, err = suite.MsgServer.CommitReport(suite.Ctx, &types.MsgCommitReport{
		Authority:  accounts[0].String(),
		ReportType: types.ReportType_REPORT_TYPE_END,
		Report:     anyVal,
	})
	suite.Require().NoError(err)
	// After ReportEnd(3)
	suite.Require().Equal(actualEmission2, suite.Keeper.GetEpochEmission(suite.Ctx, epoch2))
	suite.Require().Equal(expectedComputingPowerSum2, suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2))

	suite.Commit()
	// BeginBlocker
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetEpochEmission(suite.Ctx, epoch2))
	suite.Require().Equal(sdk.ZeroDec(), suite.Keeper.GetGlobalComputingPowerOnEpoch(suite.Ctx, epoch2))
}
