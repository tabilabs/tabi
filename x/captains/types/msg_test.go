package types

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/stretchr/testify/suite"
)

type MsgsTestSuite struct {
	suite.Suite
}

func TestMsgsTestSuite(t *testing.T) {
	suite.Run(t, new(MsgsTestSuite))
}

func (suite *MsgsTestSuite) TestMsgUpdateParamsValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgUpdateParams
		expPass   bool
	}{
		{
			"fail - invalid authority address",
			&MsgUpdateParams{
				Authority: "invalid",
				Params:    DefaultParams(),
			},
			false,
		},
		{
			"fail - no authority members",
			&MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    DefaultParams(),
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.msgUpdate.ValidateBasic()
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}

func (suite *MsgsTestSuite) TestMsgMintValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgCreateCaptainNode
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgCreateCaptainNode{
				Authority:  "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Owner:      "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				DivisionId: "1",
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgCreateCaptainNode{
				Authority:  "invalid",
				Owner:      "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				DivisionId: "1",
			},
			false,
		},
		{
			"fail - invalid Receiver address",
			&MsgCreateCaptainNode{
				Authority:  "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Owner:      "invalid",
				DivisionId: "1",
			},
			false,
		},
		{
			"fail - invalid DivisionId",
			&MsgCreateCaptainNode{
				Authority:  "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Owner:      "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				DivisionId: "",
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.msgUpdate.ValidateBasic()
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}

}

// TestMsgUpdateNodeInfoValidateBasic tests MsgUpdateNodeInfo ValidateBasic
func (suite *MsgsTestSuite) TestMsgClaimComputingPowerValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgClaimComputingPower
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgClaimComputingPower{
				Sender:               "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				ComputingPowerAmount: 100,
				NodeId:               "1",
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgClaimComputingPower{
				Sender:               "invalid",
				ComputingPowerAmount: 100,
				NodeId:               "1",
			},
			false,
		},
		{
			"fail - invalid Experience",
			&MsgClaimComputingPower{
				Sender:               "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				ComputingPowerAmount: 0,
				NodeId:               "1",
			},
			false,
		},
		{
			"fail - invalid NodeId",
			&MsgClaimComputingPower{
				Sender:               "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				ComputingPowerAmount: 100,
				NodeId:               "",
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.msgUpdate.ValidateBasic()
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}

func (suite *MsgsTestSuite) TestMsgCommitComputingPowerValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgCommitComputingPower
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgCommitComputingPower{
				Authority: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				ComputingPowerRewards: []ClaimableComputingPower{
					{
						Owner:  "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
						Amount: 100,
					},
				},
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgCommitComputingPower{
				Authority: "invalid",
				ComputingPowerRewards: []ClaimableComputingPower{
					{
						Owner:  "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
						Amount: 100,
					},
				},
			},
			false,
		},
		{
			"fail - invalid Receiver address",
			&MsgCommitComputingPower{
				Authority: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				ComputingPowerRewards: []ClaimableComputingPower{
					{
						Owner:  "invalid",
						Amount: 100,
					},
				},
			},
			false,
		},
		{
			"fail - invalid Experience",
			&MsgCommitComputingPower{
				Authority: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				ComputingPowerRewards: []ClaimableComputingPower{
					{
						Owner:  "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
						Amount: 0,
					},
				},
			},
			false,
		},
		{
			"fail - invalid Experience",
			&MsgCommitComputingPower{
				Authority:             "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				ComputingPowerRewards: nil,
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.msgUpdate.ValidateBasic()
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}

func (suite *MsgsTestSuite) TestMsgAddAuthorizedMembersValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgAddAuthorizedMembers
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgAddAuthorizedMembers{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Members: []string{
					"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				},
			},
			true,
		},
		{
			"fail - invalid authority address",
			&MsgAddAuthorizedMembers{
				Authority: "invalid",
				Members: []string{
					"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				},
			},
			false,
		},
		{
			"fail - invalid caller address",
			&MsgAddAuthorizedMembers{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Members: []string{
					"invalid",
				},
			},
			false,
		},
		{
			"fail - invalid caller is empty",
			&MsgAddAuthorizedMembers{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Members:   []string{},
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.msgUpdate.ValidateBasic()
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}

func (suite *MsgsTestSuite) TestMsgRemoveAuthorizedMembersValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgRemoveAuthorizedMembers
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgRemoveAuthorizedMembers{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Members: []string{
					"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				},
			},
			true,
		},
		{
			"fail - invalid authority address",
			&MsgRemoveAuthorizedMembers{
				Authority: "invalid",
				Members: []string{
					"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				},
			},
			false,
		},
		{
			"fail - invalid caller address",
			&MsgRemoveAuthorizedMembers{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Members: []string{
					"invalid",
				},
			},
			false,
		},
		{
			"fail - invalid caller is empty",
			&MsgRemoveAuthorizedMembers{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Members:   []string{},
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.msgUpdate.ValidateBasic()
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}

func (suite *MsgsTestSuite) TestMsgUpdateSaleLevelValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgUpdateSaleLevel
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgUpdateSaleLevel{
				Authority: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				SaleLevel: 2,
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgUpdateSaleLevel{
				Authority: "invalid",
				SaleLevel: 2,
			},
			false,
		},
		{
			"fail - invalid SaleLevel",
			&MsgUpdateSaleLevel{
				Authority: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				SaleLevel: 0,
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.msgUpdate.ValidateBasic()
			if tc.expPass {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}
