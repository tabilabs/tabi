package keeper

import (
	"context"

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

func (m msgServer) RegisterCaller(
	goCtx context.Context,
	msg *types.MsgRegisterCaller,
) (*types.MsgRegisterCallerResponse, error) {
	if m.k.authority.String() != msg.Authority {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid authority; expected %s, got %s",
			m.k.authority.String(),
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	m.k.RegisterCallers(ctx, msg.Callers)
	return &types.MsgRegisterCallerResponse{}, nil
}

/*****************************************************************************/

// Mint implement the interface of types.MsgServer
func (m msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, msg.Sender) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"invalid sender; not in allow list",
		)
	}

	// generate node id
	nodeId := m.k.GenerateNodeID(ctx)
	node := types.NewNode(nodeId, msg.DivisionId, msg.Receiver)
	if err := m.k.CreateNode(ctx, node, sender); err != nil {
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

func (m msgServer) ReceiveExperience(goCtx context.Context, msg *types.MsgReceiveExperience) (*types.MsgReceiveExperienceResponse, error) {
	// todo
	return &types.MsgReceiveExperienceResponse{}, nil
}

func (m msgServer) UpdatePowerOnPeriod(goCtx context.Context, msg *types.MsgUpdatePowerOnPeriod) (*types.MsgUpdatePowerOnPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// check if msg.Sender not in allow list
	if !m.k.AuthCaller(ctx, msg.Sender) {
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
