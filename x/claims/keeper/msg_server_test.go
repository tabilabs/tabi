package keeper_test

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/tabilabs/tabi/x/claims/types"
)

func (suite *ClaimsTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr bool
	}{
		{
			name:      "fail - invalid authority",
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
		suite.Run("MsgUpdateParams", func() {
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
	testCases := []struct {
		name      string
		request   *types.MsgClaims
		expectErr bool
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
		},
	}
	for _, tc := range testCases {
		suite.mockCk.SetCurrentEpoch(10)
		suite.Run("MsgClaims", func() {
			_, err := suite.msgServer.Claims(suite.ctx, tc.request)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
