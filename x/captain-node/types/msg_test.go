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
			"pass - valid msg",
			&MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    DefaultParams(),
			},
			true,
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
		msgUpdate *MsgMint
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgMint{
				Sender:     "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Receiver:   "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				DivisionId: "1",
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgMint{
				Sender:     "invalid",
				Receiver:   "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				DivisionId: "1",
			},
			false,
		},
		{
			"fail - invalid Receiver address",
			&MsgMint{
				Sender:     "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Receiver:   "invalid",
				DivisionId: "1",
			},
			false,
		},
		{
			"fail - invalid DivisionId",
			&MsgMint{
				Sender:     "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Receiver:   "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
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

func (suite *MsgsTestSuite) TestMsgWithdrawExperienceValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgWithdrawExperience
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgWithdrawExperience{
				Sender:     "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Experience: 100,
				NodeId:     "1",
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgWithdrawExperience{
				Sender:     "invalid",
				Experience: 100,
				NodeId:     "1",
			},
			false,
		},
		{
			"fail - invalid Experience",
			&MsgWithdrawExperience{
				Sender:     "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Experience: 0,
				NodeId:     "1",
			},
			false,
		},
		{
			"fail - invalid NodeId",
			&MsgWithdrawExperience{
				Sender:     "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				Experience: 100,
				NodeId:     "",
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

func (suite *MsgsTestSuite) TestMsgUpdatePowerOnPeriodValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgUpdatePowerOnPeriod
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgUpdatePowerOnPeriod{
				Sender: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				CaptainNodePowerOnPeriods: []*CaptainNodePowerOnPeriod{
					{
						NodeId:        "1",
						PowerOnPeriod: 18,
					},
				},
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgUpdatePowerOnPeriod{
				Sender: "invalid",
				CaptainNodePowerOnPeriods: []*CaptainNodePowerOnPeriod{
					{
						NodeId:        "1",
						PowerOnPeriod: 18,
					},
				},
			},
			false,
		},
		{
			"fail - invalid NodeId",
			&MsgUpdatePowerOnPeriod{
				Sender: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				CaptainNodePowerOnPeriods: []*CaptainNodePowerOnPeriod{
					{
						NodeId:        "",
						PowerOnPeriod: 18,
					},
				},
			},
			false,
		},
		{
			"fail - invalid PowerOnPeriod",
			&MsgUpdatePowerOnPeriod{
				Sender: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				CaptainNodePowerOnPeriods: []*CaptainNodePowerOnPeriod{
					{
						NodeId:        "1",
						PowerOnPeriod: 0,
					},
				},
			},
			false,
		},
		{
			"fail - invalid PowerOnPeriod",
			&MsgUpdatePowerOnPeriod{
				Sender:                    "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				CaptainNodePowerOnPeriods: nil,
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

func (suite *MsgsTestSuite) TestMsgUpdateUserExperienceValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgUpdateUserExperience
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgUpdateUserExperience{
				Sender: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				UserExperiences: []*UserExperience{
					{
						Receiver:   "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
						Experience: 100,
					},
				},
			},
			true,
		},
		{
			"fail - invalid Sender address",
			&MsgUpdateUserExperience{
				Sender: "invalid",
				UserExperiences: []*UserExperience{
					{
						Receiver:   "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
						Experience: 100,
					},
				},
			},
			false,
		},
		{
			"fail - invalid Receiver address",
			&MsgUpdateUserExperience{
				Sender: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				UserExperiences: []*UserExperience{
					{
						Receiver:   "invalid",
						Experience: 100,
					},
				},
			},
			false,
		},
		{
			"fail - invalid Experience",
			&MsgUpdateUserExperience{
				Sender: "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				UserExperiences: []*UserExperience{
					{
						Receiver:   "cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
						Experience: 0,
					},
				},
			},
			false,
		},
		{
			"fail - invalid Experience",
			&MsgUpdateUserExperience{
				Sender:          "cosmos1vewsdxxmeraett7ztsaym88jsrv85kzm8ekjsg",
				UserExperiences: nil,
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

func (suite *MsgsTestSuite) TestMsgRegisterCallerValidateBasic() {
	testCases := []struct {
		name      string
		msgUpdate *MsgRegisterCaller
		expPass   bool
	}{
		{
			"pass - valid msg",
			&MsgRegisterCaller{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Callers: []string{
					"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				},
			},
			true,
		},
		{
			"fail - invalid authority address",
			&MsgRegisterCaller{
				Authority: "invalid",
				Callers: []string{
					"cosmos1pgm8hyk0pvphmlvfjc8wsvk4daluz5tgmr4lac",
				},
			},
			false,
		},
		{
			"fail - invalid caller address",
			&MsgRegisterCaller{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Callers: []string{
					"invalid",
				},
			},
			false,
		},
		{
			"fail - invalid caller is empty",
			&MsgRegisterCaller{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Callers:   []string{},
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
