package keeper_test

import sdk "github.com/cosmos/cosmos-sdk/types"

func (suite *IntegrationTestSuite) TestEpochState() {
	es := NewEpochState(suite)
	es.InitNodes(suite, accounts[0].String(), 1, 100)
	es.InitNodesPowerOnRatio()

	crp := NewCaptainsReporter(sdk.OneDec(), 10)

	for i := 0; i < 30; i++ {
		es = ExecuteEpoch(es, crp)
	}
}
