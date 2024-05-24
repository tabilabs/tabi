package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/version"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/spf13/cobra"

	"github.com/tabilabs/tabi/x/claims/types"
)

// GetQueryCmd returns the cli query commands for the mint module.
func GetQueryCmd() *cobra.Command {
	cliamsQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the cliams module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cliamsQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryRewards(),
	)
	return cliamsQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current minting parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryRewards implements a command to return the owner's rewards.
func GetCmdQueryRewards() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewards [owner]",
		Short: "Query the owner's rewards",
		Long: fmt.Sprintf(`Query the owner's rewards

Example:
$ %s query %s rewards <owner_addr>
`, version.AppName, types.ModuleName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.HolderTotalRewards(
				context.Background(),
				&types.QueryHolderTotalRewardsRequest{Owner: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
