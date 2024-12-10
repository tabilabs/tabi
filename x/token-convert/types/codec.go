package types

import (
	cryptocodec "github.com/tabilabs/tabi/crypto/codec"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

const (
	// Amino names
	convertTabiName   = "tokenconvert/MsgConvertTabi"
	convertVetabiName = "tokenconvert/MsgConvertVetabi"

	cancelConvertName = "tokenconvert/MsgCancelConvert"

	withdrawName = "tokenconvert/MsgWithdrawTabi"
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgConvertTabi{},
		&MsgConvertVetabi{},
		&MsgWithdrawTabi{},
		&MsgCancelConvert{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgConvertTabi{}, convertTabiName, nil)
	cdc.RegisterConcrete(&MsgConvertVetabi{}, convertVetabiName, nil)
	cdc.RegisterConcrete(&MsgCancelConvert{}, cancelConvertName, nil)
	cdc.RegisterConcrete(&MsgWithdrawTabi{}, withdrawName, nil)
}
