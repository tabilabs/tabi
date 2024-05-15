package keeper_test

import (
	"errors"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/tabilabs/tabi/testutil"
	tabitypes "github.com/tabilabs/tabi/types"
	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *CaptainsTestSuite) utilsAddAuthorizedMember(member string) {
	suite.msgServer.AddAuthorizedMembers(
		suite.ctx,
		&types.MsgAddAuthorizedMembers{
			Authority: accounts[0].String(),
			Members:   []string{member},
		})
}

func (suite *CaptainsTestSuite) utilsCreateCaptainNode(owner string, divisionLevel uint64) string {
	divisions := suite.utilsGetDivisions()
	resp, _ := suite.msgServer.CreateCaptainNode(
		suite.ctx,
		&types.MsgCreateCaptainNode{
			Authority:  accounts[0].String(),
			Owner:      owner,
			DivisionId: divisions[divisionLevel],
		},
	)
	return resp.NodeId
}

func (suite *CaptainsTestSuite) utilsBatchCreateCaptainNode(owner string, divisionLevel, amount uint64) []string {
	nodeIds := make([]string, amount)
	for i := uint64(0); i < amount; i++ {
		nodeIds[i] = suite.utilsCreateCaptainNode(owner, divisionLevel)
	}
	return nodeIds
}

func (suite *CaptainsTestSuite) utilsBatchAssignFixedPowerOnRatio(nodes []string, value int64, prec int64) []types.NodePowerOnRatio {
	nodePowers := make([]types.NodePowerOnRatio, len(nodes))
	for i, node := range nodes {
		nodePowers[i] = types.NodePowerOnRatio{
			NodeId:           node,
			OnOperationRatio: sdk.NewDecWithPrec(value, prec),
		}
	}
	return nodePowers
}

func (suite *CaptainsTestSuite) utilsBatchAssignRandomPowerOnRatio(nodes []string) []types.NodePowerOnRatio {
	rand.Seed(suite.ctx.BlockTime().Unix())
	nodePowers := make([]types.NodePowerOnRatio, len(nodes))
	for i, node := range nodes {
		power, _ := sdk.NewDecFromStr(fmt.Sprintf("%f", 0.47+rand.Float64()*0.53))
		nodePowers[i] = types.NodePowerOnRatio{
			NodeId:           node,
			OnOperationRatio: power,
		}
	}
	return nodePowers
}

func (suite *CaptainsTestSuite) utilsGetDivisions() map[uint64]string {
	resp, err := suite.queryClient.Divisions(suite.ctx, &types.QueryDivisionsRequest{})
	suite.NoError(err)

	divisionMap := make(map[uint64]string)
	for _, division := range resp.Divisions {
		divisionMap[division.Level] = division.Id
	}

	return divisionMap
}

func (suite *CaptainsTestSuite) utilsCommitPower(owner string, amount uint64) {
	_, err := suite.msgServer.CommitComputingPower(suite.ctx, &types.MsgCommitComputingPower{
		Authority: accounts[0].String(),
		ComputingPowerRewards: []types.ClaimableComputingPower{
			{amount, owner},
		},
	})
	suite.NoError(err)
}

func (suite *CaptainsTestSuite) utilsUpdateLevel(level uint64) {
	_, err := suite.msgServer.UpdateSaleLevel(suite.ctx, &types.MsgUpdateSaleLevel{
		Authority: accounts[0].String(),
		SaleLevel: level,
	})
	suite.NoError(err)
}

// utilsStakingTabiWithAmount delegates amount of atabi to a validator.
func (suite *CaptainsTestSuite) utilsStakingTabiWithAmount(
	delegator sdk.AccAddress,
	amount int64,
	validator stakingtypes.Validator,
) {
	_, err := suite.app.StakingKeeper.Delegate(
		suite.ctx,
		delegator,
		sdk.NewInt(amount),
		1,
		validator,
		true)
	suite.Require().NoError(err)
}

func (suite *CaptainsTestSuite) utilsFundToken(addr sdk.AccAddress, amt int64, denom string) error {
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
