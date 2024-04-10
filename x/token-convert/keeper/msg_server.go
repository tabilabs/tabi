package keeper

import (
	"context"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) ConvertTabi(goCtx context.Context, msg *types.MsgConvertTabi) (*types.MsgConvertTabiResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) ConvertVetabi(goCtx context.Context, msg *types.MsgConvertVetabi) (*types.MsgConvertVetabiResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) WithdrawTabi(goCtx context.Context, msg *types.MsgWithdrawTabi) (*types.MsgWithdrawTabiResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) CancelConvert(goCtx context.Context, msg *types.MsgCancelConvert) (*types.MsgCancelConvertResponse, error) {
	//TODO implement me
	panic("implement me")
}
