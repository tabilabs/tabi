package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tabilabs/tabi/x/claims/types"
)

type msgServer struct {
	k Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the mint MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

// UpdateParams implement the interface of types.MsgServer
func (m msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if m.k.authority.String() != msg.Authority {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority.String(),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}
	return &types.MsgUpdateParamsResponse{}, nil
}

// WithdrawNodeReward implement the interface of types.MsgServer
func (m msgServer) Claims(goCtx context.Context, msg *types.MsgClaims) (*types.MsgClaimsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// if receiver is empty, set receiver to sender
	if len(msg.Receiver) == 0 {
		msg.Receiver = msg.Sender
	}
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	// check if the sender has not held node
	if !m.k.HasNode(ctx, sender) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrUnauthorized,
			"sender %s has not held node",
			msg.Sender,
		)
	}

	amount, err := m.k.WithdrawRewards(ctx, sender, receiver)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimsResponse{
		Amount: amount,
	}, nil
}
