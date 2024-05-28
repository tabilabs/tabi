package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tabilabs/tabi/x/limiter/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the limiter MsgServer interface
func NewMsgServerImpl(keeper *Keeper) msgServer {
	return msgServer{keeper}
}

// UpdateParams defines a method that allows to update the parameters of the module
// NOTE: use x/params instead before sdk v0.47.
func (m msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if m.k.authority.String() != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"invalid authority: expected %s, got %s", m.k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetParamsInModule(ctx, msg.Params); err != nil {
		return nil, errorsmod.Wrap(err, "failed to set params")
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
