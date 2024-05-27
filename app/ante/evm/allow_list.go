package evm

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/ethereum/go-ethereum/common"

	evmtypes "github.com/tabilabs/tabi/x/evm/types"
)

// EthAllowListVerificationDecorator validates an account is allowed to send a transaction
type EthAllowListVerificationDecorator struct {
	ak            evmtypes.AccountKeeper
	evmKeeper     EVMKeeper
	limiterKeeper LimiterKeeper
}

// NewEthAllowListVerificationDecorator creates a new EthAllowListVerificationDecorator
func NewEthAllowListVerificationDecorator(ak evmtypes.AccountKeeper, ek EVMKeeper, lk LimiterKeeper) EthAllowListVerificationDecorator {
	return EthAllowListVerificationDecorator{
		ak:            ak,
		evmKeeper:     ek,
		limiterKeeper: lk,
	}
}

func (alvd EthAllowListVerificationDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	if !ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	// skip if limiter is disabled
	if !alvd.limiterKeeper.IsEnabled(ctx) {
		return next(ctx, tx, simulate)
	}

	for _, msg := range tx.GetMsgs() {
		msgEthTx, ok := msg.(*evmtypes.MsgEthereumTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrUnknownRequest, "invalid message type %T, expected %T", msg, (*evmtypes.MsgEthereumTx)(nil))
		}

		// skip if the transaction is not a contract creation
		to := msgEthTx.AsTransaction().To()
		if to != nil {
			continue
		}

		from := msgEthTx.GetFrom()
		acc := alvd.ak.GetAccount(ctx, from)
		if acc == nil {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrUnknownAddress,
				"account %s is nil", common.BytesToAddress(msgEthTx.GetFrom().Bytes()),
			)
		}

		if !alvd.limiterKeeper.IsAuthorized(ctx, acc.GetAddress()) {
			return ctx, errorsmod.Wrapf(
				errortypes.ErrUnauthorized,
				"account %s is not allowed to send transactions", common.BytesToAddress(msgEthTx.GetFrom().Bytes()),
			)
		}
	}
	return next(ctx, tx, simulate)
}
