package keeper_test

import sdk "github.com/cosmos/cosmos-sdk/types"

func (suite *CaptainsTestSuite) TestComputingPower() {
	addr1 := accounts[1].String()
	node1 := suite.utilsCreateCaptainNode(addr1, 1)

	// epoch 1: 2000 * e ^ (20/3) = 1571543.988
	epoch1 := suite.keeper.GetCurrentEpoch(suite.ctx)
	power1 := suite.keeper.CalcNodeComputingPowerOnEpoch(
		suite.ctx,
		epoch1,
		node1,
		sdk.NewDecWithPrec(1, 0),
	)
	suite.Require().Equal(sdk.NewDecWithPrec(1571543988, 3), power1)

	// epoch2: 2000 * e ^ (1/2) = 5436.564
	suite.keeper.EnterNewEpoch(suite.ctx)
	epoch2 := suite.keeper.GetCurrentEpoch(suite.ctx)
	power2 := suite.keeper.CalcNodeComputingPowerOnEpoch(
		suite.ctx,
		epoch2,
		node1,
		sdk.NewDecWithPrec(1, 0),
	)
	suite.Require().Equal(sdk.NewDecWithPrec(5606325, 2), power2)
}
