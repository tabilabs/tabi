package keeper_test

import (
	tabitypes "github.com/tabilabs/tabi/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestSampleOwnerPledge() {
	owner := accounts[1]
	err := suite.utilsFundToken(owner, 10_000_000, tabitypes.AttoTabi)
	suite.Require().NoError(err)

	// staking tokens
	ownerAddr := sdk.MustAccAddressFromBech32(owner.String())
	for i := 0; i < len(suite.Validators); i++ {
		suite.utilsStakingTabiWithAmount(ownerAddr, 1_000_000, suite.Validators[i])
	}

	val, err := suite.Keeper.SampleOwnerPledge(suite.Ctx, ownerAddr)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewDec(10_000_000), val)
	suite.T().Logf("SampleOwnerPledge: %s", val.String())
}

func (suite *IntegrationTestSuite) TestCalcGlobalPledgeRatio() {
	testCases := []struct {
		name      string
		melleate  func()
		expectVal sdk.Dec
	}{
		{
			name:      "claimed sum: 0, pledge sum: 0",
			melleate:  func() {},
			expectVal: sdk.OneDec(),
		},
		{
			name: "claimed sum: 100, pledge sum: 0",
			melleate: func() {
				suite.Keeper.SetGlobalClaimedEmission(suite.Ctx, sdk.MustNewDecFromStr("100"))
			},
			expectVal: sdk.MustNewDecFromStr("0.3"),
		},
		{
			name: "claimed sum: 100, pledge sum: 1",
			melleate: func() {
				suite.Keeper.SetGlobalClaimedEmission(suite.Ctx, sdk.MustNewDecFromStr("100"))
				suite.Keeper.SetGlobalPledge(suite.Ctx, 1, sdk.MustNewDecFromStr("1"))
			},
			expectVal: sdk.MustNewDecFromStr("0.3"),
		},
		{
			name: "claimed sum: 100, pledge sum: 50",
			melleate: func() {
				suite.Keeper.SetGlobalClaimedEmission(suite.Ctx, sdk.MustNewDecFromStr("100"))
				suite.Keeper.SetGlobalPledge(suite.Ctx, 1, sdk.MustNewDecFromStr("50"))
			},
			expectVal: sdk.MustNewDecFromStr("0.5"),
		},
		{
			name: "claimed sum: 100, pledge sum: 120",
			melleate: func() {
				suite.Keeper.SetGlobalClaimedEmission(suite.Ctx, sdk.MustNewDecFromStr("100"))
				suite.Keeper.SetGlobalPledge(suite.Ctx, 1, sdk.MustNewDecFromStr("120"))
			},
			expectVal: sdk.MustNewDecFromStr("1.0"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.melleate()

			val := suite.Keeper.CalcGlobalPledgeRatio(suite.Ctx, 1)
			suite.Require().Equal(tc.expectVal, val)
		})
	}
}

func (suite *IntegrationTestSuite) TestCalcNodePledgeRatioOnEpoch() {
	owner := accounts[1]
	nodeId := suite.utilsCreateCaptainNode(owner.String(), 1)

	testCases := []struct {
		name      string
		melleate  func()
		expectVal sdk.Dec
	}{
		{
			name:      "claimed sum: 0, pledge sum: 0",
			melleate:  func() {},
			expectVal: sdk.MustNewDecFromStr("0.3"),
		},
		{
			name: "claimed sum: 100, pledge sum: 0",
			melleate: func() {
				suite.Keeper.SetNodeClaimedEmission(suite.Ctx, nodeId, sdk.MustNewDecFromStr("100"))
				suite.Keeper.SetOwnerPledge(suite.Ctx, owner, 1, sdk.MustNewDecFromStr("0"))
			},
			expectVal: sdk.MustNewDecFromStr("0"),
		},
		{
			name: "claimed sum: 100, pledge sum: 1",
			melleate: func() {
				suite.Keeper.SetNodeClaimedEmission(suite.Ctx, nodeId, sdk.MustNewDecFromStr("100"))
				suite.Keeper.SetOwnerPledge(suite.Ctx, owner, 1, sdk.MustNewDecFromStr("1"))
			},
			expectVal: sdk.MustNewDecFromStr("0.01"),
		},
		{
			name: "claimed sum: 100, pledge sum: 50",
			melleate: func() {
				suite.Keeper.SetNodeClaimedEmission(suite.Ctx, nodeId, sdk.MustNewDecFromStr("100"))
				suite.Keeper.SetOwnerPledge(suite.Ctx, owner, 1, sdk.MustNewDecFromStr("50"))
			},
			expectVal: sdk.MustNewDecFromStr("0.3"),
		},
		{
			name: "claimed sum: 100, pledge sum: 120",
			melleate: func() {
				suite.Keeper.SetNodeClaimedEmission(suite.Ctx, nodeId, sdk.MustNewDecFromStr("100"))
				suite.Keeper.SetOwnerPledge(suite.Ctx, owner, 1, sdk.MustNewDecFromStr("120"))
			},
			expectVal: sdk.MustNewDecFromStr("0.3"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.melleate()

			val := suite.Keeper.CalcNodePledgeRatioOnEpoch(suite.Ctx, 1, nodeId)
			suite.Require().Equal(tc.expectVal, val)
		})
	}
}
