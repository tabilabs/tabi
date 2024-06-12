package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tabiante "github.com/tabilabs/tabi/app/ante/tabi"
	"github.com/tabilabs/tabi/testutil"
	tabitypes "github.com/tabilabs/tabi/types"
	captainstypes "github.com/tabilabs/tabi/x/captains/types"
	claimstypes "github.com/tabilabs/tabi/x/claims/types"
)

//  TestEpochState tests the epoch state transitions as follows:
//  TODO: revise this state table.
//
//	+--------------------+------------------------------+---------------------------------------+
//	| Phase              | State                        | Expectation                           |
//	|--------------------+------------------------------+---------------------------------------|
//	| BeginEpoch         | EpochEmission(T-1)           | Check if deleted                      |
//	|                    | GlobalComputingPower(T-1)    |                                       |
//	|                    | ReportDigest(T), EndEpoch(T) |                                       |
//	|                    | Batches(T), StandByFlag      |                                       |
//	|--------------------+------------------------------+---------------------------------------|
//	| BeforeReportDigest | GlobalPledge(T)              | 1. If T == 1, skip                    |
//	|                    |                              | 2. Check existence                    |
//	|                    | GlobalClaimedEmission        | 1. If T == 1, skip                    |
//	|                    |                              | 2. If none claimed, skip              |
//	|                    |                              | 3. If claimed, check existence and    |
//	|                    |                              |    equality                           |
//	|                    | ExpectEpochEmission          | Calculate directly                    |
//	|--------------------+------------------------------+---------------------------------------|
//	| AfterReportDigest  | StandByOverFlag              | Check existence                       |
//	|                    | GlobalPledge(T)              | 1. If T == 1, skip                    |
//	|                    |                              | 2. Check non-existence                |
//	|                    | ActualEpochEmission          | Read from storage, and ensure it      |
//	|                    |                              | equals ExpectEpochEmission            |
//	|--------------------+------------------------------+---------------------------------------|
//	| BeforeReportBatch  | NodeCumulativeEmission       | 1. If T == 1, skip                    |
//	|                    | ByEpoch(T-1)                 | 2. If T == 2, calculate               |
//	|                    |                              |    NodeCumulativeEmissionByEpoch(T-1) |
//	|                    |                              | 3. If T >= 3, calculate               |
//	|                    |                              |    NodeCumulativeEmissionByEpoch(T-1) |
//	|                    |                              |    and check NodeCumulativeEmission   |
//	|                    |                              |    ByEpoch(T-2) existence             |
//	|                    | OwnerPledge(T)               | 1. If T == 1, skip?                   |
//	|                    |                              | 2. If T >= 2, check existence of      |
//	|                    |                              |    OwnerPledge(T)                     |
//	|                    | ExpectComputingPower(N,T)    | Calculate ExpectComputingPower        |
//	|--------------------+------------------------------+---------------------------------------|
//	| AfterReportBatch   | NodeCumulativeEmission       | 1. If T == 1, skip?                   |
//	|                    | ByEpoch(T-1)                 | 2. Read                               |
//	|                    |                              |    NodeCumulativeEmissionByEpoch(T-1) |
//	|                    |                              | 3. Check non-existence of             |
//	|                    |                              |    NodeCumulativeEmissionByEpoch(T-2) |
//	|                    | OwnerPledge(T)               | 1. If T > 2, check deletion of        |
//	|                    |                              |    OwnerPledge(T-2)                   |
//	|                    | ActualComputingPower(N,T)    | Read from storage, ensure it equals   |
//	|                    |                              |    ExpectComputingPower               |
//	+--------------------+------------------------------+---------------------------------------+

func (suite *IntegrationTestSuite) TestEpochState() {
	owner := accounts[0].String()

	// todo: separate cases into [epoch-phase-a, epoch-phase-b] under the same context.
	testCases := []EpochTestCase{
		{
			name:      "no staking and claiming",
			maxEpoch:  3,
			reporter:  NewCaptainsReporter(sdk.OneDec(), 10),
			saveState: true,
		},
		{
			name:     "claiming",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			execStandByFn: func(es *EpochState) {
				if es.Epoch >= 2 {
					_, err := suite.ClaimsServer.Claims(suite.Ctx, &claimstypes.MsgClaims{
						Receiver: owner,
						Sender:   owner,
					})
					suite.Require().NoError(err)
				}
			},
			saveState: true,
		},
		{
			name:     "claiming and staking",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			execStandByFn: func(es *EpochState) {
				if es.Epoch == 1 {
					// staking
					err := suite.utilsFundToken(sdk.MustAccAddressFromBech32(owner), 10_000_000, tabitypes.AttoTabi)
					suite.Require().NoError(err)
					// staking tokens
					ownerAddr := sdk.MustAccAddressFromBech32(owner)
					for i := 0; i < len(suite.Validators); i++ {
						suite.utilsStakingTabiWithAmount(ownerAddr, 1_000_000, suite.Validators[i])
					}
				}

				if es.Epoch >= 2 {
					_, err := suite.ClaimsServer.Claims(suite.Ctx, &claimstypes.MsgClaims{
						Receiver: owner,
						Sender:   owner,
					})
					suite.Require().NoError(err)
				}
			},
			saveState: true,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			// NOTE: don't set this in the test case as for each test the suite will evaluate it before the test thus
			// creating too much nodes.
			tc.currEpochState = NewEpochState(suite, tc.reporter).WithNodes(owner, 1, 100).WithNodesPowerOnRatio("")
			tc.WithStateMap()
			tc.Execute()
		})
		suite.SetupTest()
	}
}

// TestEpochBusy tests ante handler restrictions on certain messages during the busy phase.
func (suite *IntegrationTestSuite) TestAnteWhenEpochBusy() {
	owner := accounts[0].String()

	testCases := []EpochTestCase{
		{
			name:     "fail to create note",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			execBusyFn: func(es *EpochState) {
				dec := tabiante.NewCaptainsRestrictionDecorator(suite.Keeper)
				err := testutil.ValidateAnteForMsgs(suite.Ctx, dec, &captainstypes.MsgCreateCaptainNode{
					Authority:  owner,
					Owner:      owner,
					DivisionId: "not-empty",
				})
				suite.Require().ErrorContains(err, "not allowed in busy phrase")
			},
		},
		{
			name:     "fail to update sale level",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			execBusyFn: func(es *EpochState) {
				dec := tabiante.NewCaptainsRestrictionDecorator(suite.Keeper)
				err := testutil.ValidateAnteForMsgs(suite.Ctx, dec, &captainstypes.MsgUpdateSaleLevel{
					Authority: owner,
					SaleLevel: 2,
				})
				suite.Require().ErrorContains(err, "not allowed in busy phrase")
			},
		},
		{
			name:     "fail to commit computing power",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			execBusyFn: func(es *EpochState) {
				dec := tabiante.NewCaptainsRestrictionDecorator(suite.Keeper)
				err := testutil.ValidateAnteForMsgs(suite.Ctx, dec, &captainstypes.MsgCommitComputingPower{
					Authority:             owner,
					ComputingPowerRewards: nil,
				})
				suite.Require().ErrorContains(err, "not allowed in busy phrase")
			},
		},
		{
			name:     "fail to claim computing power",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			execBusyFn: func(es *EpochState) {
				dec := tabiante.NewCaptainsRestrictionDecorator(suite.Keeper)
				err := testutil.ValidateAnteForMsgs(suite.Ctx, dec, &captainstypes.MsgClaimComputingPower{
					Sender:               owner,
					ComputingPowerAmount: 100,
					NodeId:               "not-empty",
				})
				suite.Require().ErrorContains(err, "not allowed in busy phrase")
			},
		},
		{
			name:     "fail to claim rewards",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			execBusyFn: func(es *EpochState) {
				dec := tabiante.NewCaptainsRestrictionDecorator(suite.Keeper)
				err := testutil.ValidateAnteForMsgs(suite.Ctx, dec, &claimstypes.MsgClaims{
					Receiver: owner,
					Sender:   owner,
				},
				)
				suite.Require().ErrorContains(err, "not allowed in busy phrase")
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.currEpochState = NewEpochState(suite, tc.reporter).WithNodes(owner, 1, 100).WithNodesPowerOnRatio("")
			tc.Execute()
		})
		suite.SetupTest()
	}
}

// TestEpochStateScenarioA tests the epoch state transitions as follows:
// global pledge ratio: ?
// owner pledge ratio: 0.3

func (suite *IntegrationTestSuite) TestEpochStateScenarioA() {
	owner := accounts[0].String()
	testCases := []EpochTestCase{
		{
			name:              "test epoch state scenario A",
			maxEpoch:          2,
			reporter:          NewCaptainsReporter(sdk.MustNewDecFromStr("0.47"), 100),
			saveState:         true,
			globalPledgeRatio: sdk.MustNewDecFromStr("0.8"),
			execStandByFn: func(es *EpochState) {
				if es.Epoch >= 2 {
					resp, err := suite.ClaimsServer.Claims(suite.Ctx, &claimstypes.MsgClaims{
						Receiver: owner,
						Sender:   owner,
					})
					suite.Require().NoError(err)
					suite.T().Logf("claimed: %s", resp.Amount.String())
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.currEpochState = NewEpochState(suite, tc.reporter).
				WithNodes(owner, 1, 100).
				WithNodesPowerOnRatio("1.0")

			tc.WithStateMap()
			tc.Execute()
		})
		suite.SetupTest()
	}
}
