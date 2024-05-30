package cli

import (
	"github.com/spf13/cobra"

	tabitypes "github.com/tabilabs/tabi/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token convert transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewTxCmdConvertTabi(),
		NewTxCmdConvertVetabi(),
		NewTxCmdWithdrawTabi(),
		NewTxCmdCancelConvert(),
	)

	return cmd
}

// NewTxCmdConvertTabi is the cli cmd for ConvertTabi
func NewTxCmdConvertTabi() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-tabi [amount]",
		Short: "Convert tabi to vetabi",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			if coin.Denom != tabitypes.AttoTabi {
				return errorsmod.Wrapf(types.ErrInvalidCoin, "invalid coin denom: %s", coin.Denom)
			}

			msg := types.NewMsgConvertTabi(coin, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewTxCmdConvertVetabi is the cli cmd for ConvertVetabi
func NewTxCmdConvertVetabi() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert-vetabi [amount] [strategy-name]",
		Short: "Convert vetabi to tabi as per allowed strategy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			if coin.Denom != tabitypes.AttoVeTabi {
				return errorsmod.Wrapf(types.ErrInvalidCoin, "invalid coin denom: %s", coin.Denom)
			}

			msg := types.NewMsgConvertVetabi(coin, clientCtx.GetFromAddress(), args[1])

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewTxCmdWithdrawTabi is the cli cmd for WithdrawTabi
func NewTxCmdWithdrawTabi() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-tabi [voucher-id]",
		Short: "Withdraw tabi as per the voucher",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawTabi(args[0], clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewTxCmdCancelConvert is the cli cmd for CancelConvert
func NewTxCmdCancelConvert() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-convert [voucher-id]",
		Short: "Cancel converting vetabi to tabi",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelConvert(args[0], clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
