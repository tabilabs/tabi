package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *IntegrationTestSuite) TestEpochState() {

	testCases := []EpochTestCase{
		{
			name:     "no staking and claiming",
			maxEpoch: 3,
			reporter: NewCaptainsReporter(sdk.OneDec(), 10),
			currEpochState: NewEpochState(suite).
				WithNodes(accounts[0].String(), 1, 100).
				WithNodesPowerOnRatio(),
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.Execute()
		})
	}
}
