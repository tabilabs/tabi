package keeper_test

import (
	"fmt"

	"github.com/tabilabs/tabi/x/claims/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *ClaimsTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr bool
	}{
		{
			name: "fail - invalid authority",

			request:   &types.MsgUpdateParams{Authority: "foobar"},
			expectErr: true,
		},
		{
			name: "pass - valid Update msg",
			request: &types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    types.DefaultParams(),
			},
			expectErr: false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgUpdateParams - %s", tc.name), func() {
			_, err := suite.msgServer.UpdateParams(suite.ctx, tc.request)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *ClaimsTestSuite) TestClaims() {
	// setup bank keeper
	testCases := []struct {
		name      string
		request   *types.MsgClaims
		expectErr bool
		setup     func() *MockCaptains
	}{
		{
			name:      "fail - invalid sender",
			request:   &types.MsgClaims{Sender: "foobar"},
			expectErr: true,
		},
		{
			name:      "fail - invalid receiver",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: "foobar"},
			expectErr: true,
		},
		{
			name:      "fail - nodes length is 0",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: true,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase01)
			},
		},
		// has 1 node start
		{
			name:      "fail - owner has 1 node, epoch == 1 the node has no rewards to claim and has no historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: true,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase02)
			},
		},
		{
			name:      "success - owner has 1 node, epoch == 2 the node has rewards to claim and has no historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: false,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase03)
			},
		},
		{
			name:      "fail - owner has 1 node, epoch == 2 the node has no rewards to claim and has historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: true,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase04)
			},
		},
		{
			name:      "success - owner has 1 node, epoch == 3 the node has rewards to claim and has historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: false,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase05)
			},
		},
		// has 1 node end
		// has 5 node start
		{
			name:      "fail - owner has 5 node, epoch == 1 the node has no rewards to claim and has no historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: true,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase06)
			},
		},
		{
			name:      "success - owner has 5 node, epoch == 2 the node has rewards to claim and has no historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: false,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase07)
			},
		},
		{
			name:      "fail - owner has 5 node, epoch == 2 the node has no rewards to claim and has historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: true,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase08)
			},
		},
		{
			name:      "success - owner has 5 node, epoch == 3 the node has rewards to claim and has historical emission",
			request:   &types.MsgClaims{Sender: suite.cosmosAddress.String(), Receiver: suite.cosmosAddress.String()},
			expectErr: false,
			setup: func() *MockCaptains {
				return NewMockCaptains(KeyCase09)
			},
		},
		// has 5 node end

	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("MsgClaims - %s", tc.name), func() {
			if tc.setup != nil {
				mockCaptains := tc.setup()
				suite.app.ClaimsKeeper.SetCaptainsKeeper(mockCaptains)
			}
			_, err := suite.msgServer.Claims(suite.ctx, tc.request)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
