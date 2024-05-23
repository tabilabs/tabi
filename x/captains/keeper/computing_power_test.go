package keeper_test

import sdk "github.com/cosmos/cosmos-sdk/types"

func (suite *IntegrationTestSuite) TestComputingPower() {
	addr1 := accounts[1].String()
	node1 := suite.utilsCreateCaptainNode(addr1, 1)

	// epoch 1: 2000 * e ^ (2) = 14778.112
	epoch1 := suite.Keeper.GetCurrentEpoch(suite.Ctx)
	power1 := suite.Keeper.CalcNodeComputingPowerOnEpoch(
		suite.Ctx,
		epoch1,
		node1,
		sdk.NewDecWithPrec(1, 0),
	)
	suite.Require().Equal(sdk.NewDecWithPrec(14778112, 3), power1)
}
