package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestValidateGenesis() {
	testCases := []struct {
		name     string
		genState *GenesisState
		expPass  bool
	}{
		{
			name:     "TestValidateGenesisWithValidParams",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			// TODO: complete test here
			name: "TestValidateGenesisWithInvalidParams",
			genState: &GenesisState{
				Params: Params{
					// Set invalid params here
				},
				Divisions:                DefaultDivision(),
				Nodes:                    nil,
				EpochState:               EpochState{},
				ClaimableComputingPowers: nil,
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.genState.Validate()
		if tc.expPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}
