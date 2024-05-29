package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tabilabs/tabi/x/limiter/types"
)

func (suite *IntegrationTestSuite) TestLimiterSwitch() {
	testCases := []struct {
		name      string
		enable    bool
		authority string
		expectErr bool
	}{
		{
			name:      "success: enable limiter",
			enable:    true,
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
		},
		{
			name:      "success: disable limiter",
			enable:    false,
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
		},
		{
			name:      "failure: unauthorized",
			enable:    true,
			authority: accounts[0].String(),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.MsgServer.LimiterSwitch(suite.Ctx, &types.MsgLimiterSwitch{
				Authority: tc.authority,
				Enabled:   tc.enable,
			})
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().Equal(tc.enable, suite.App.LimiterKeeper.IsEnabled(suite.Ctx))
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestAddAllowListMember() {
	testCases := []struct {
		name      string
		melleate  func()
		authority string
		member    string
		expectErr bool
	}{
		{
			name:      "success: add member",
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
			member:    accounts[0].String(),
		},
		{
			name:      "failure: member already exists",
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
			melleate: func() {
				suite.App.LimiterKeeper.AddAllowListMember(suite.Ctx, accounts[0].String())
			},
			member:    accounts[0].String(),
			expectErr: true,
		},
		{
			name:      "failure: unauthorized",
			authority: accounts[0].String(),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.melleate != nil {
				tc.melleate()
			}

			_, err := suite.MsgServer.AddAllowListMember(suite.Ctx, &types.MsgAddAllowListMember{
				Authority: tc.authority,
				Address:   tc.member,
			})
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().True(suite.App.LimiterKeeper.IsAuthorized(suite.Ctx,
					sdk.MustAccAddressFromBech32(tc.member)))
			}

			suite.SetupTest() // reset
		})
	}
}

func (suite *IntegrationTestSuite) TestRemoveAllowListMember() {

	testCases := []struct {
		name      string
		melleate  func()
		authority string
		member    string
		expectErr bool
	}{
		{
			name:      "success: remove member",
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
			melleate: func() {
				suite.App.LimiterKeeper.AddAllowListMember(suite.Ctx, accounts[0].String())
			},
			member: accounts[0].String(),
		},
		{
			name:      "success: remove member",
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
			melleate: func() {
				suite.App.LimiterKeeper.AddAllowListMember(suite.Ctx, accounts[0].String())
				suite.App.LimiterKeeper.AddAllowListMember(suite.Ctx, accounts[1].String())
				suite.App.LimiterKeeper.AddAllowListMember(suite.Ctx, accounts[2].String())
			},
			member: accounts[2].String(),
		},
		{
			name:      "failure: member not found",
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
			member:    accounts[0].String(),
			melleate: func() {
				suite.App.LimiterKeeper.AddAllowListMember(suite.Ctx, accounts[1].String())
			},
			expectErr: true,
		},
		{
			name:      "failure: allow list is empty",
			authority: "tabis10d07y265gmmuvt4z0w9aw880jnsr700j7ry74f",
			member:    accounts[0].String(),
			expectErr: true,
		},
		{
			name:      "failure: unauthorized",
			authority: accounts[0].String(),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			if tc.melleate != nil {
				tc.melleate()
			}

			_, err := suite.MsgServer.RemoveAllowListMember(suite.Ctx, &types.MsgRemoveAllowListMember{
				Authority: tc.authority,
				Address:   tc.member,
			})

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().False(suite.App.LimiterKeeper.IsAuthorized(suite.Ctx,
					sdk.MustAccAddressFromBech32(tc.member)))
			}

			suite.SetupTest() // reset
		})
	}
}
