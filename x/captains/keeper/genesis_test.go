package keeper_test

import (
	"github.com/tabilabs/tabi/x/captains/types"
)

func (suite *IntegrationTestSuite) TestExportAndImportGenesis() {

	testCases := []struct {
		name      string
		prepare   func()
		expectErr bool
	}{
		{
			name: "default",
			prepare: func() {
				suite.afterEpochOne()
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prepare()

			nodes, _ := suite.QueryClient.Nodes(suite.Ctx, &types.QueryNodesRequest{})

			// export and validate
			gs1 := suite.Keeper.ExportGenesis(suite.Ctx)
			err := gs1.Validate()
			suite.Require().NoError(err)

			suite.Require().Equal(nodes.Nodes[0].Id, gs1.NodesComputingPower[0].NodeId)

			// reset the state
			suite.utilsPruneCaptainsStore()

			// import and export again
			suite.Keeper.InitGenesis(suite.Ctx, gs1)
			gs2 := suite.Keeper.ExportGenesis(suite.Ctx)

			suite.Require().Equal(gs1.Params, gs2.Params)
			suite.Require().Equal(gs1.BaseState, gs2.BaseState)
			suite.Require().Equal(gs1.Divisions, gs2.Divisions)
			suite.Require().Equal(gs1.Nodes, gs2.Nodes)
			suite.Require().Equal(gs1.EpochesEmission, gs2.EpochesEmission)
			suite.Require().Equal(gs1.NodesClaimedEmission, gs2.NodesClaimedEmission)
			suite.Require().Equal(gs1.NodesCumulativeEmission, gs2.NodesCumulativeEmission)
			suite.Require().Equal(gs1.GlobalsPledge, gs2.GlobalsPledge)
			suite.Require().Equal(gs1.OwnersPledge, gs2.OwnersPledge)
			suite.Require().Equal(gs1.OwnersClaimableComputingPower, gs2.OwnersClaimableComputingPower)
			suite.Require().Equal(gs1.GlobalsComputingPower, gs2.GlobalsComputingPower)
			suite.Require().Equal(gs1.Batches, gs2.Batches)

			suite.Require().Equal(gs1.NodesComputingPower, gs2.NodesComputingPower)
		})
	}

}

func (suite *IntegrationTestSuite) utilsPruneCaptainsStore() {
	key := suite.App.GetKey(types.StoreKey)
	store := suite.Ctx.KVStore(key)

	iter := store.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

func (suite *IntegrationTestSuite) CompareGenesis(gs1, gs2 *types.GenesisState) {
	suite.Require().Equal(gs1.BaseState, gs2.BaseState)
	suite.Require().Equal(gs1.Params, gs2.Params)
	suite.Require().Equal(gs1.Divisions, gs2.Divisions)
	suite.Require().Equal(gs1.Nodes, gs2.Nodes)
	suite.Require().Equal(gs1.EpochesEmission, gs2.EpochesEmission)
	suite.Require().Equal(gs1.NodesClaimedEmission, gs2.NodesClaimedEmission)
	suite.Require().Equal(gs1.NodesCumulativeEmission, gs2.NodesCumulativeEmission)
	suite.Require().Equal(gs1.GlobalsPledge, gs2.GlobalsPledge)
	suite.Require().Equal(gs1.OwnersPledge, gs2.OwnersPledge)
	suite.Require().Equal(gs1.OwnersClaimableComputingPower, gs2.OwnersClaimableComputingPower)
}
