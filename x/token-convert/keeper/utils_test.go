package keeper_test

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/testutil"
	tabitypes "github.com/tabilabs/tabi/types"
)

func (suite *TokenConvertTestSuite) utilsFundToken(addr sdk.AccAddress, amt int64, denom string) error {
	coins := make([]sdk.Coin, 1)

	switch denom {
	case tabitypes.AttoTabi:
		coins[0] = tabitypes.NewTabiCoinInt64(amt)
	case tabitypes.AttoVeTabi:
		coins[0] = tabitypes.NewVeTabiCoinInt64(amt)
	default:
		return errors.New("unsupported denom")
	}

	return testutil.FundAccount(suite.ctx, suite.app.BankKeeper, addr, sdk.NewCoins(coins...))
}
