package keeper

import (
	"context"

	"github.com/tabilabs/tabi/x/limiter/types"
)

type MsgServer struct {
	*Keeper
}

func NewMsgServerImpl(keeper *Keeper) MsgServer {
	return MsgServer{keeper}
}

// UpdateParams defines a method that allows to update the parameters of the module
// NOTE: use x/params instead before sdk v0.47.
func (m MsgServer) UpdateParams(ctx context.Context, params *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}
