package types

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisTestSuite struct {
	suite.Suite
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestValidateGenesis() {
	testCases := []struct {
		name      string
		genState  *GenesisState
		expectErr bool
	}{
		{
			name:      "success: default genesis state",
			genState:  DefaultGenesisState(),
			expectErr: false,
		},
		{
			name: "fail: base state",
			genState: &GenesisState{
				Params: DefaultParams(),
				BaseState: BaseState{
					EpochId:               0,
					NextNodeSequence:      0,
					GlobalClaimedEmission: sdk.ZeroDec(),
				},
				Divisions: DefaultDivision(),
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.genState.Validate()
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
