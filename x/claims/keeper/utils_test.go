package keeper_test

import (
	"errors"

	utiltx "github.com/tabilabs/tabi/testutil/tx"

	"github.com/tabilabs/tabi/testutil"
	tabitypes "github.com/tabilabs/tabi/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	KeyCase01 = iota + 1
	KeyCase02
	KeyCase03
	KeyCase04
	KeyCase05
	KeyCase06
	KeyCase07
	KeyCase08
	KeyCase09
)

const (
	KeyQueryNodeTotalRewards01 = iota + 100
	KeyQueryNodeTotalRewards02
	KeyQueryNodeTotalRewards03
	KeyQueryNodeTotalRewards04
)

const (
	KeyQueryHolderTotalRewards01 = iota + 200
	KeyQueryHolderTotalRewards02
	KeyQueryHolderTotalRewards03
	KeyQueryHolderTotalRewards04
	KeyQueryHolderTotalRewards05
	KeyQueryHolderTotalRewards06
	KeyQueryHolderTotalRewards07
)

var accounts = []sdk.AccAddress{
	sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
	sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
	sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
}

func (suite *ClaimsTestSuite) utilsFundToken(addr sdk.AccAddress, amt int64, denom string) error {
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
