package keeper_test

import sdk "github.com/cosmos/cosmos-sdk/types"

func (suite *CaptainsTestSuite) TestComputingPower() {
	addr1 := accounts[1].String()
	node1 := suite.utilsCreateCaptainNode(addr1, 1)

	// epoch 1: 2000 * e ^ (2) = 14778.112
	epoch1 := suite.keeper.GetCurrentEpoch(suite.ctx)
	power1 := suite.keeper.CalcNodeComputingPowerOnEpoch(
		suite.ctx,
		epoch1,
		node1,
		sdk.NewDecWithPrec(1, 0),
	)
	suite.Require().Equal(sdk.NewDecWithPrec(14778112, 3), power1)
}
