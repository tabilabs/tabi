package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tabilabs/tabi/x/token-convert/types"
)

// NewQueryCmd returns the cli query commands for this module
func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Token convert query subcommands",
		DisableFlagParsing: true,
	}

	cmd.AddCommand(
		NewQueryCmdStrategy(),
		NewQueryCmdStrategies(),
		NewQueryCmdVoucher(),
		NewQueryCmdVouchers(),
		NewQueryCmdVoucherStatus(),
	)

	return cmd
}

// NewQueryCmdStrategy is the cli cmd for QueryStrategy
func NewQueryCmdStrategy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategy [name]",
		Short: "Query a strategy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Strategy(
				context.Background(),
				&types.QueryStrategyRequest{
					Name: args[0],
				})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryCmdStrategies is the cli cmd for QueryStrategies
func NewQueryCmdStrategies() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategies",
		Short: "Query strategies",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Strategies(
				context.Background(),
				&types.QueryStrategiesRequest{
					Pagination: pageReq,
				})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "strategies")

	return cmd
}

// NewQueryCmdVoucher is the cli cmd for QueryVoucher
func NewQueryCmdVoucher() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voucher [id]",
		Short: "Query a voucher",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Voucher(
				context.Background(),
				&types.QueryVoucherRequest{
					VoucherId: args[0],
				})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// NewQueryCmdVouchers is the cli cmd for QueryVouchers
func NewQueryCmdVouchers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vouchers",
		Short: "Query vouchers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			_, err = sdk.AccAddressFromBech32(owner)
			if err != nil {
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Vouchers(
				context.Background(),
				&types.QueryVouchersRequest{
					Owner:      owner,
					Pagination: pageReq,
				})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	cmd.Flags().AddFlagSet(FlagSetVouchers)
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "vouchers")

	return cmd
}

// NewQueryCmdVoucherStatus is the cli cmd for QueryVoucherStatus
func NewQueryCmdVoucherStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voucher-status [id]",
		Short: "Query a voucher status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.VoucherStatus(
				context.Background(),
				&types.QueryVoucherStatusRequest{
					VoucherId: args[0],
				})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
