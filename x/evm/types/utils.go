package types

import (
	"errors"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/gogo/protobuf/proto"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

// DefaultPriorityReduction is the default amount of price values required for 1 unit of priority.
// Because priority is `int64` while price is `big.Int`, it's necessary to scale down the range to keep it more pratical.
// The default value is the same as the `sdk.DefaultPowerReduction`.
var DefaultPriorityReduction = sdk.DefaultPowerReduction

var EmptyCodeHash = crypto.Keccak256(nil)

const maxBitLen = 256

// DecodeTxResponse decodes an protobuf-encoded byte slice into TxResponse
func DecodeTxResponse(in []byte) (*MsgEthereumTxResponse, error) {
	var txMsgData sdk.TxMsgData
	if err := proto.Unmarshal(in, &txMsgData); err != nil {
		return nil, err
	}

	if len(txMsgData.MsgResponses) == 0 {
		return &MsgEthereumTxResponse{}, nil
	}

	var res MsgEthereumTxResponse
	if err := proto.Unmarshal(txMsgData.MsgResponses[0].Value, &res); err != nil {
		return nil, errorsmod.Wrap(err, "failed to unmarshal tx response message data")
	}

	return &res, nil
}

// EncodeTransactionLogs encodes TransactionLogs slice into a protobuf-encoded byte slice.
func EncodeTransactionLogs(res *TransactionLogs) ([]byte, error) {
	return proto.Marshal(res)
}

// DecodeTransactionLogs decodes an protobuf-encoded byte slice into TransactionLogs
func DecodeTransactionLogs(data []byte) (TransactionLogs, error) {
	var logs TransactionLogs
	err := proto.Unmarshal(data, &logs)
	if err != nil {
		return TransactionLogs{}, err
	}
	return logs, nil
}

// UnwrapEthereumMsg extract MsgEthereumTx from wrapping sdk.Tx
func UnwrapEthereumMsg(tx *sdk.Tx, ethHash common.Hash) (*MsgEthereumTx, error) {
	if tx == nil {
		return nil, fmt.Errorf("invalid tx: nil")
	}

	for _, msg := range (*tx).GetMsgs() {
		ethMsg, ok := msg.(*MsgEthereumTx)
		if !ok {
			return nil, fmt.Errorf("invalid tx type: %T", tx)
		}
		txHash := ethMsg.AsTransaction().Hash()
		ethMsg.Hash = txHash.Hex()
		if txHash == ethHash {
			return ethMsg, nil
		}
	}

	return nil, fmt.Errorf("eth tx not found: %s", ethHash)
}

// BinSearch execute the binary search and hone in on an executable gas limit
func BinSearch(lo, hi uint64, executable func(uint64) (bool, *MsgEthereumTxResponse, error)) (uint64, error) {
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := executable(mid)
		// If the error is not nil(consensus error), it means the provided message
		// call or transaction will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more.
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi, nil
}

// EffectiveGasPrice compute the effective gas price based on eip-1159 rules
// `effectiveGasPrice = min(baseFee + tipCap, feeCap)`
func EffectiveGasPrice(baseFee, feeCap, tipCap *big.Int) *big.Int {
	return math.BigMin(new(big.Int).Add(tipCap, baseFee), feeCap)
}

func ValidateEthTx(tx *ethtypes.Transaction) error {
	if !IsValidInt256(tx.Value()) {
		return errors.New("value overflow")
	}
	if !IsValidInt256(tx.GasPrice()) {
		return errors.New("gas price overflow")
	}
	if !IsValidInt256(tx.GasFeeCap()) {
		return errors.New("gas fee cap overflow")
	}
	if !IsValidInt256(tx.GasTipCap()) {
		return errors.New("gas tip cap overflow")
	}
	return nil
}

// IsValidInt256 check the bound of 256 bit number
func IsValidInt256(i *big.Int) bool {
	return i == nil || i.BitLen() <= maxBitLen
}
