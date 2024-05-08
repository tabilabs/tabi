package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tabitypes "github.com/tabilabs/tabi/types"
	"github.com/tabilabs/tabi/x/token-convert/types"
)

func (suite *TokenConvertTestSuite) TestConvertTabi() {
	sender := accounts[0].String()

	suite.utilsFundToken(accounts[0], 1_000_000, tabitypes.AttoTabi)
	res := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoTabi)
	suite.Require().Equal(tabitypes.NewTabiCoinInt64(1_000_000), res)

	testCases := []struct {
		name         string
		msg          *types.MsgConvertTabi
		expectVetabi sdk.Coin
		expectTabi   sdk.Coin
		expectErr    bool
	}{
		{
			name: "success - enough balance",
			msg: &types.MsgConvertTabi{
				Coin:   tabitypes.NewTabiCoinInt64(1_000),
				Sender: sender,
			},
			expectVetabi: tabitypes.NewVeTabiCoinInt64(1_000),
			expectTabi:   tabitypes.NewTabiCoinInt64(999_000),
			expectErr:    false,
		},
		{
			name: "fail - not enough balance",
			msg: &types.MsgConvertTabi{
				Coin:   tabitypes.NewTabiCoinInt64(10_000_000),
				Sender: sender,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.msgServer.ConvertTabi(suite.ctx, tc.msg)

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tabiCoin := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoTabi)
				suite.Require().Equal(tc.expectTabi, tabiCoin)
				vetabiCoin := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoVeTabi)
				suite.Require().Equal(tc.expectVetabi, vetabiCoin)
			}

			suite.SetupTest() // reset balances
		})
	}
}

func (suite *TokenConvertTestSuite) TestConvertVetabi() {
	sender := accounts[0].String()

	testCases := []struct {
		name          string
		msg           *types.MsgConvertVetabi
		expectVetabi  sdk.Coin
		expectTabi    sdk.Coin
		expectVoucher bool
		expectErr     bool
	}{
		{
			name: "success - select instant strategy with 1,000 avetabi ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(1_000),
				Strategy: types.StrategyInstant,
				Sender:   sender,
			},
			expectVetabi:  tabitypes.NewVeTabiCoinInt64(999_000),
			expectTabi:    tabitypes.NewTabiCoinInt64(250),
			expectVoucher: false,
			expectErr:     false,
		},
		{
			name: "success - select instant strategy with 100 avetabi ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(100),
				Strategy: types.StrategyInstant,
				Sender:   sender,
			},
			expectVetabi:  tabitypes.NewVeTabiCoinInt64(999_900),
			expectTabi:    tabitypes.NewTabiCoinInt64(25),
			expectVoucher: false,
			expectErr:     false,
		},
		{
			name: "success - select instant strategy with 10 avetabi ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(10),
				Strategy: types.StrategyInstant,
				Sender:   sender,
			},
			expectVetabi:  tabitypes.NewVeTabiCoinInt64(999_990),
			expectTabi:    tabitypes.NewTabiCoinInt64(2), // 2.5
			expectVoucher: false,
			expectErr:     false,
		},
		{
			name: "success - select instant strategy with 1 avetabi ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(1),
				Strategy: types.StrategyInstant,
				Sender:   sender,
			},
			expectVetabi:  tabitypes.NewVeTabiCoinInt64(999_999),
			expectTabi:    tabitypes.NewTabiCoinInt64(0), // 0.25
			expectVoucher: false,
			expectErr:     false,
		},
		{
			name: "fail - select instant strategy but no enough balance ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(10_000_000),
				Strategy: types.StrategyInstant,
				Sender:   sender,
			},
			expectVoucher: false,
			expectErr:     true,
		},
		{
			name: "success - select 90days strategy with 10,000 avetabi ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(10_000),
				Strategy: types.Strategy90Days,
				Sender:   sender,
			},
			expectVetabi:  tabitypes.NewVeTabiCoinInt64(990_000),
			expectTabi:    tabitypes.NewTabiCoinInt64(0),
			expectVoucher: true,
			expectErr:     false,
		},
		{
			name: "fail - select 90days strategy not enough balance ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(10_000_000),
				Strategy: types.Strategy90Days,
				Sender:   sender,
			},
			expectVoucher: false,
			expectErr:     true,
		},
		{
			name: "success - select 180days strategy with 10,000 avetabi ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(10_000),
				Strategy: types.Strategy180Days,
				Sender:   sender,
			},
			expectVetabi:  tabitypes.NewVeTabiCoinInt64(990_000),
			expectTabi:    tabitypes.NewTabiCoinInt64(0),
			expectVoucher: true,
			expectErr:     false,
		},
		{
			name: "fail - select 180days strategy not enough balance ",
			msg: &types.MsgConvertVetabi{
				Coin:     tabitypes.NewVeTabiCoinInt64(10_000_000),
				Strategy: types.Strategy180Days,
				Sender:   sender,
			},
			expectVoucher: false,
			expectErr:     true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset balances

			suite.utilsFundToken(accounts[0], 1_000_000, tabitypes.AttoVeTabi)
			res := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoVeTabi)
			suite.Require().Equal(tabitypes.NewVeTabiCoinInt64(1_000_000), res)

			voucher, err := suite.msgServer.ConvertVetabi(suite.ctx, tc.msg)

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tabiCoin := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoTabi)
				suite.Require().Equal(tc.expectTabi, tabiCoin)
				vetabiCoin := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoVeTabi)
				suite.Require().Equal(tc.expectVetabi, vetabiCoin)

				if tc.expectVoucher {
					resp, err := suite.queryClient.Voucher(suite.ctx, &types.QueryVoucherRequest{
						VoucherId: voucher.VoucherId,
					})
					suite.Require().NoError(err)
					suite.Require().Equal(voucher.VoucherId, resp.Id)
				}
			}
		})
	}
}

func (suite *TokenConvertTestSuite) TestWithdrawTabi() {
	sender := accounts[0].String()

	testCases := []struct {
		name         string
		msg          *types.MsgWithdrawTabi
		melleate     func(*types.MsgWithdrawTabi)
		timeAfter    time.Duration
		expectTabi   sdk.Coin
		expectVetabi sdk.Coin
		expectErr    bool
	}{
		{
			name: "success - strategy 90 days, withdraw after 10 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy90Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			timeAfter:    10 * 24 * time.Hour,
			expectTabi:   tabitypes.NewTabiCoinInt64(55_555),    // truncate(round(10^6 * 10 / 90) * 0.5)
			expectVetabi: tabitypes.NewVeTabiCoinInt64(888_889), // round(10^6 * 80 / 90) == 888,889
			expectErr:    false,
		},
		{
			name: "success - strategy 90 days, withdraw after 45 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy90Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			timeAfter:    45 * 24 * time.Hour,
			expectTabi:   tabitypes.NewTabiCoinInt64(250_000),   // truncate(round(10^6 * 45 / 90) * 0.5)
			expectVetabi: tabitypes.NewVeTabiCoinInt64(500_000), // round(10^6 * 45 / 90)
			expectErr:    false,
		},
		{
			name: "success - strategy 90 days, withdraw after 90 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy90Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			expectTabi:   tabitypes.NewTabiCoinInt64(500_000), // truncate(round(10^6 * 90 / 90) * 0.5)
			expectVetabi: tabitypes.NewVeTabiCoinInt64(0),     // round(10^6 * 0 / 90)
			timeAfter:    90 * 24 * time.Hour,
			expectErr:    false,
		},
		{
			name: "success - strategy 90 days, withdraw after 100 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy90Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			expectTabi:   tabitypes.NewTabiCoinInt64(500_000), // truncate(round(10^6 * 90 / 90) * 0.5)
			expectVetabi: tabitypes.NewVeTabiCoinInt64(0),     // round(10^6 * 0 / 90)
			timeAfter:    100 * 24 * time.Hour,
			expectErr:    false,
		},
		{
			name: "fail - strategy 90 days, withdraw immediately",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy90Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			timeAfter: 0,
			expectErr: true,
		},
		{
			name: "success - strategy 180 days, withdraw after 10 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy180Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			timeAfter:    10 * 24 * time.Hour,
			expectTabi:   tabitypes.NewTabiCoinInt64(55_556),    // round(10^6 * 10 / 180)
			expectVetabi: tabitypes.NewVeTabiCoinInt64(944_444), // round(10^6 * 170 / 180)
			expectErr:    false,
		},
		{
			name: "success - strategy 180 days, withdraw after 40 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy180Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			timeAfter:    40 * 24 * time.Hour,
			expectTabi:   tabitypes.NewTabiCoinInt64(222_222),   // truncate(10^6 * 40 / 180)
			expectVetabi: tabitypes.NewVeTabiCoinInt64(777_778), // truncate(10^6 * 140 / 180)
			expectErr:    false,
		},
		{
			name: "success - strategy 180 days, withdraw after 180 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy180Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			timeAfter:    180 * 24 * time.Hour,
			expectTabi:   tabitypes.NewTabiCoinInt64(1_000_000),
			expectVetabi: tabitypes.NewVeTabiCoinInt64(0),
			expectErr:    false,
		},
		{
			name: "success - strategy 180 days, withdraw after 200 days",
			msg: &types.MsgWithdrawTabi{
				Sender: sender,
			},
			melleate: func(msg *types.MsgWithdrawTabi) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy180Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			timeAfter:    200 * 24 * time.Hour,
			expectTabi:   tabitypes.NewTabiCoinInt64(1_000_000),
			expectVetabi: tabitypes.NewVeTabiCoinInt64(0),
			expectErr:    false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			// fund token
			suite.utilsFundToken(accounts[0], 1_000_000, tabitypes.AttoVeTabi)
			res := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoVeTabi)
			suite.Require().Equal(tabitypes.NewVeTabiCoinInt64(1_000_000), res)

			// convert and wait a period
			tc.melleate(tc.msg)
			suite.CommitAfter(tc.timeAfter)

			// withdraw
			resp, err := suite.msgServer.WithdrawTabi(suite.ctx, tc.msg)
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tabiCoin := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoTabi)
				suite.Require().Equal(tc.expectTabi, tabiCoin)
				vetabiCoin := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoVeTabi)
				suite.Require().Equal(tc.expectVetabi, vetabiCoin)

				suite.Require().Equal(resp.TabiWithdrawn, tc.expectTabi)
				suite.Require().Equal(resp.VetabiReturned, tc.expectVetabi)
			}
		})
	}
}

func (suite *TokenConvertTestSuite) TestCancelConvert() {
	sender := accounts[0].String()
	another := accounts[1].String()

	testCases := []struct {
		name         string
		melleate     func(*types.MsgCancelConvert)
		msg          *types.MsgCancelConvert
		timeAfter    time.Duration
		expectVetabi sdk.Coin
		expectErr    bool
	}{
		{
			name: "success - create but cancel later",
			msg: &types.MsgCancelConvert{
				Sender: sender,
			},
			melleate: func(msg *types.MsgCancelConvert) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy90Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			expectVetabi: tabitypes.NewVeTabiCoinInt64(1_000_000),
			expectErr:    false,
		},
		{
			name: "fail - not voucher owner",
			msg: &types.MsgCancelConvert{
				Sender: another,
			},
			melleate: func(msg *types.MsgCancelConvert) {
				resp, err := suite.msgServer.ConvertVetabi(suite.ctx, &types.MsgConvertVetabi{
					Coin:     tabitypes.NewVeTabiCoinInt64(1_000_000),
					Strategy: types.Strategy90Days,
					Sender:   sender,
				})
				suite.Require().NoError(err)
				msg.VoucherId = resp.VoucherId
			},
			expectErr: true,
		},
		{
			name: "fail - voucher not found",
			msg: &types.MsgCancelConvert{
				Sender: sender,
			},
			melleate: func(msg *types.MsgCancelConvert) {
				msg.VoucherId = "voucher-not-found"
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			// fund token
			suite.utilsFundToken(accounts[0], 1_000_000, tabitypes.AttoVeTabi)
			res := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoVeTabi)
			suite.Require().Equal(tabitypes.NewVeTabiCoinInt64(1_000_000), res)

			// convert and wait a period
			tc.melleate(tc.msg)

			_, err := suite.msgServer.CancelConvert(suite.ctx, tc.msg)

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				vetabiCoin := suite.bankKeeper.GetBalance(suite.ctx, accounts[0], tabitypes.AttoVeTabi)
				suite.Require().Equal(tc.expectVetabi, vetabiCoin)
			}
		})
	}
}
