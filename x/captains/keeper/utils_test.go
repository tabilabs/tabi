package keeper_test

import (
	"errors"
	"fmt"
	"math/rand"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tabilabs/tabi/testutil"
	tabitypes "github.com/tabilabs/tabi/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *IntegrationTestSuite) utilsAddAuthorizedMember(member string) {
	suite.MsgServer.AddAuthorizedMembers(
		suite.Ctx,
		&types.MsgAddAuthorizedMembers{
			Authority: accounts[0].String(),
			Members:   []string{member},
		})
}

func (suite *IntegrationTestSuite) utilsCreateCaptainNode(owner string, divisionLevel uint64) string {
	divisions := suite.utilsGetDivisions()
	resp, _ := suite.MsgServer.CreateCaptainNode(
		suite.Ctx,
		&types.MsgCreateCaptainNode{
			Authority:  accounts[0].String(),
			Owner:      owner,
			DivisionId: divisions[divisionLevel],
		},
	)
	return resp.NodeId
}

func (suite *IntegrationTestSuite) utilsBatchCreateCaptainNode(owner string, divisionLevel, amount uint64) []string {
	nodeIds := make([]string, amount)
	for i := uint64(0); i < amount; i++ {
		nodeIds[i] = suite.utilsCreateCaptainNode(owner, divisionLevel)
	}
	return nodeIds
}

func (suite *IntegrationTestSuite) utilsBatchAssignFixedPowerOnRatio(nodes []string, value int64, prec int64) []types.NodePowerOnRatio {
	nodePowers := make([]types.NodePowerOnRatio, len(nodes))
	for i, node := range nodes {
		nodePowers[i] = types.NodePowerOnRatio{
			NodeId:           node,
			OnOperationRatio: sdk.NewDecWithPrec(value, prec),
		}
	}
	return nodePowers
}

func (suite *IntegrationTestSuite) utilsBatchAssignRandomPowerOnRatio(nodes []string) []types.NodePowerOnRatio {
	rand.Seed(suite.Ctx.BlockTime().Unix())
	nodePowers := make([]types.NodePowerOnRatio, len(nodes))
	for i, node := range nodes {
		power, _ := sdk.NewDecFromStr(fmt.Sprintf("%f", 0.47+rand.Float64()*0.53))
		nodePowers[i] = types.NodePowerOnRatio{
			NodeId:           node,
			OnOperationRatio: power,
		}
	}
	return nodePowers
}

func (suite *IntegrationTestSuite) utilsGetDivisions() map[uint64]string {
	resp, err := suite.QueryClient.Divisions(suite.Ctx, &types.QueryDivisionsRequest{})
	suite.NoError(err)

	divisionMap := make(map[uint64]string)
	for _, division := range resp.Divisions {
		divisionMap[division.Level] = division.Id
	}

	return divisionMap
}

func (suite *IntegrationTestSuite) utilsCommitPower(owner string, amount uint64) {
	_, err := suite.MsgServer.CommitComputingPower(suite.Ctx, &types.MsgCommitComputingPower{
		Authority: accounts[0].String(),
		ComputingPowerRewards: []types.ClaimableComputingPower{
			{amount, owner},
		},
	})
	suite.NoError(err)
}

func (suite *IntegrationTestSuite) utilsUpdateLevel(level uint64) {
	_, err := suite.MsgServer.UpdateSaleLevel(suite.Ctx, &types.MsgUpdateSaleLevel{
		Authority: accounts[0].String(),
		SaleLevel: level,
	})
	suite.NoError(err)
}

// utilsStakingTabiWithAmount delegates amount of atabi to a validator.
func (suite *IntegrationTestSuite) utilsStakingTabiWithAmount(
	delegator sdk.AccAddress,
	amount int64,
	validator stakingtypes.Validator,
) {
	_, err := suite.App.StakingKeeper.Delegate(
		suite.Ctx,
		delegator,
		sdk.NewInt(amount),
		1,
		validator,
		true)
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) utilsFundToken(addr sdk.AccAddress, amt int64, denom string) error {
	coins := make([]sdk.Coin, 1)

	switch denom {
	case tabitypes.AttoTabi:
		coins[0] = tabitypes.NewTabiCoinInt64(amt)
	case tabitypes.AttoVeTabi:
		coins[0] = tabitypes.NewVeTabiCoinInt64(amt)
	default:
		return errors.New("unsupported denom")
	}

	return testutil.FundAccount(suite.Ctx, suite.App.BankKeeper, addr, sdk.NewCoins(coins...))
}

func (suite *IntegrationTestSuite) afterEpochOne() {
	// create 10 nodes for addr1
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
}
