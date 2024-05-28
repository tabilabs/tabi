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
func (m msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", msg.Authority)
	}

	if m.k.authority.String() != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"invalid authority: expected %s, got %s", m.k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetParams(ctx, msg.Params); err != nil {
		return nil, errorsmod.Wrap(err, "failed to set params")
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// LimiterSwitch defines a method that allows to enable or disable the limiter
func (m msgServer) LimiterSwitch(goCtx context.Context, msg *types.MsgLimiterSwitch) (*types.MsgLimiterSwitchResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", msg.Authority)
	}

	if m.k.authority.String() != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"invalid authority: expected %s, got %s", m.k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.SetEnabled(ctx, msg.Enabled); err != nil {
		return nil, errorsmod.Wrap(err, "failed to set limiter switch")
	}

	return &types.MsgLimiterSwitchResponse{}, nil
}

// AddAllowListMember defines a method that allows to add a member to the allow list
func (m msgServer) AddAllowListMember(goCtx context.Context, msg *types.MsgAddAllowListMember) (*types.MsgAddAllowListMemberResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", msg.Authority)
	}

	if m.k.authority.String() != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"invalid authority: expected %s, got %s", m.k.authority, msg.Authority)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", msg.Address)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.AddAllowListMember(ctx, msg.Address); err != nil {
		return nil, errorsmod.Wrap(err, "failed to add allow list member")
	}

	return &types.MsgAddAllowListMemberResponse{}, nil
}

// RemoveAllowListMember defines a method that allows to remove a member from the allow list
func (m msgServer) RemoveAllowListMember(goCtx context.Context, msg *types.MsgRemoveAllowListMember) (*types.MsgRemoveAllowListMemberResponse, error) {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", msg.Authority)
	}

	if m.k.authority.String() != msg.Authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"invalid authority: expected %s, got %s", m.k.authority, msg.Authority)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Address); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address: %s", msg.Address)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.k.RemoveAllowListMember(ctx, msg.Address); err != nil {
		return nil, errorsmod.Wrap(err, "failed to remove allow list member")
	}

	return &types.MsgRemoveAllowListMemberResponse{}, nil
}
