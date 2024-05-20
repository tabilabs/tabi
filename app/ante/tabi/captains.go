package tabi

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	anteutlis "github.com/tabilabs/tabi/app/ante/utils"
	captainstypes "github.com/tabilabs/tabi/x/captains/types"
	claimestypes "github.com/tabilabs/tabi/x/claims/types"
)

// CaptainsRestrictionDecorator restricts certain messages from being processed
type CaptainsRestrictionDecorator struct {
	captainsKeeper anteutlis.CaptainsKeeper
}

// NewCaptainsRestrictionDecorator creates a new CaptainsRestrictionDecorator
func NewCaptainsRestrictionDecorator(keeper anteutlis.CaptainsKeeper) CaptainsRestrictionDecorator {
	return CaptainsRestrictionDecorator{
		captainsKeeper: keeper,
	}
}

// AnteHandle restricts certain messages from being processed
func (cld CaptainsRestrictionDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	if err := cld.restrict(ctx, tx.GetMsgs()); err != nil {
		return ctx, errorsmod.Wrapf(errortypes.ErrInvalidRequest, err.Error())
	}

	return next(ctx, tx, simulate)
}

// restrict checks if the given messages are restricted
func (cld CaptainsRestrictionDecorator) restrict(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, msg := range msgs {
		switch msg := msg.(type) {
		case *captainstypes.MsgUpdateParams,
			*captainstypes.MsgCreateCaptainNode,
			*captainstypes.MsgUpdateSaleLevel,
			*captainstypes.MsgClaimComputingPower,
			*captainstypes.MsgCommitComputingPower,
			*claimestypes.MsgClaims:
			if !cld.captainsKeeper.IsStandByPhase(ctx) {
				return fmt.Errorf("msg %s is not allowed in busy phrase", msg.String())
			}
		case
			*captainstypes.MsgCommitReport:
			if cld.captainsKeeper.IsStandByPhase(ctx) {
				return fmt.Errorf("msg %s is not allowed in stand-by phrase", msg.String())
			}
		}
	}
	return nil
}
