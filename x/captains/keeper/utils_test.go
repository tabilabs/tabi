package keeper_test

import (
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
