package keeper_test

func (suite *IntegrationTestSuite) TestEpochState() {
	// prepare nods
	es := NewEpochState(suite)
	es.InitNodes(suite, accounts[0].String(), 1, 100)
	es.InitNodesPowerOnRatio()

	// check
	// TODO: start epoch testing
}
