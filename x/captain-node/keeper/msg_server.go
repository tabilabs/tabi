package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tabilabs/tabi/x/captain-node/types"
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
func (m msgServer) UpdateParams(
	goCtx context.Context,
	msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
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

/*****************************************************************************/
/*****************************************************************************/
/* Need Allow Function */
/*****************************************************************************/
/*****************************************************************************/

// Mint implement the interface of types.MsgServer
func (m msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in allow list",
		)
	}

	// generate node id
	nodeId := m.k.GenerateNodeID(ctx)
	node := types.NewNode(nodeId, msg.DivisionId, msg.Receiver)
	if err := m.k.CreateNode(ctx, node, receiver); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNode,
			sdk.NewAttribute(types.AttributeKeyNodeID, nodeId),
			sdk.NewAttribute(types.AttributeKeyDivisionID, msg.DivisionId),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgMintResponse{}, nil
}

func (m msgServer) UpdatePowerOnPeriod(goCtx context.Context, msg *types.MsgUpdatePowerOnPeriod) (*types.MsgUpdatePowerOnPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in allow list",
		)
	}

	if len(msg.CaptainNodePowerOnPeriods) == 0 {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid period; empty",
		)
	}

	events := m.k.UpdateAllNodesPowerOnPeriod(ctx, msg.CaptainNodePowerOnPeriods)

	resultEvents := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	)}
	resultEvents = append(resultEvents, events...)
	ctx.EventManager().EmitEvents(resultEvents)
	return &types.MsgUpdatePowerOnPeriodResponse{}, nil
}

func (m msgServer) UpdateUserExperience(goCtx context.Context, msg *types.MsgUpdateUserExperience) (*types.MsgUpdateUserExperienceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in allow list",
		)
	}
	if len(msg.UserExperiences) == 0 {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"invalid user experience; empty",
		)
	}
	events := m.k.UpdateAllUserExperience(ctx, msg.UserExperiences)

	resultEvents := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	)}
	resultEvents = append(resultEvents, events...)
	ctx.EventManager().EmitEvents(resultEvents)
	return &types.MsgUpdateUserExperienceResponse{}, nil
}

func (m msgServer) UpdateSaleLevel(goCtx context.Context, msg *types.MsgUpdateSaleLevel) (*types.MsgUpdateSaleLevelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in allow list",
		)
	}

	event, err := m.k.UpdateSaleLevel(ctx, msg.SaleLevel)
	if err != nil {
		return nil, err
	}
	resultEvents := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	)}
	resultEvents = append(resultEvents, event)
	ctx.EventManager().EmitEvents(resultEvents)

	return &types.MsgUpdateSaleLevelResponse{}, nil
}

func (m msgServer) AddCaller(goCtx context.Context, msg *types.MsgAddCaller) (*types.MsgAddCallerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in allow list",
		)
	}

	events, err := m.k.SetCaller(ctx, msg.Callers)
	if err != nil {
		return nil, err
	}

	resultEvents := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	)}
	resultEvents = append(resultEvents, events...)
	ctx.EventManager().EmitEvents(resultEvents)

	return &types.MsgAddCallerResponse{}, nil
}

func (m msgServer) RemoveCaller(goCtx context.Context, msg *types.MsgRemoveCaller) (*types.MsgRemoveCallerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in allow list",
		)
	}

	events, err := m.k.RemoveCaller(ctx, msg.Callers)
	if err != nil {
		return nil, err
	}

	resultEvents := sdk.Events{sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	)}
	resultEvents = append(resultEvents, events...)
	ctx.EventManager().EmitEvents(resultEvents)

	return &types.MsgRemoveCallerResponse{}, nil
}

/*****************************************************************************/
/*****************************************************************************/
/* User Function */
/*****************************************************************************/
/*****************************************************************************/

func (m msgServer) WithdrawExperience(goCtx context.Context, msg *types.MsgWithdrawExperience) (*types.MsgWithdrawExperienceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	if err := m.k.UpdateNode(ctx, msg.NodeId, msg.Experience, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNode,
			sdk.NewAttribute(types.AttributeKeyNodeID, msg.NodeId),
			sdk.NewAttribute(types.AttributeKeyExperience, fmt.Sprintf("%d", msg.Experience)),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgWithdrawExperienceResponse{}, nil
}
