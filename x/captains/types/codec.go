package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCaptainNode{},
		&MsgCommitReport{},
		&MsgAddAuthorizedMembers{},
		&MsgRemoveAuthorizedMembers{},
		&MsgUpdateSaleLevel{},
		&MsgCommitComputingPower{},
		&MsgClaimComputingPower{},
		&MsgUpdateParams{},
	)

	registry.RegisterImplementations((*ReportContent)(nil),
		&ReportDigest{},
		&ReportBatch{},
		&ReportEmission{},
		&ReportEnd{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
