package cli

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/spf13/cobra"

	"github.com/tabilabs/tabi/x/captains/types"
)

// GetQueryCmd returns the cli query commands for the captains module.
func GetQueryCmd() *cobra.Command {
	captionNodeQueryCmd := &cobra.Command{
		Use:                        "captain-node",
		Short:                      "Querying commands for the captains module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	captionNodeQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetSupplyCmd(),
		GetDivisionCmd(),
		GetDivisionsCmd(),
		GetNodeCmd(),
		GetNodesCmd(),
		GetSaleLevelCmd(),
		GetAuthorizedMembersCmd(),
	)
	return captionNodeQueryCmd
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

// GetNodeCmd returns the command to query a node
func GetNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node [node-id]",
		Short: "Query the node details",
		Long: fmt.Sprintf(`Query the division details

Example:
$ %s query %s node <node-id>
`, version.AppName, types.ModuleName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Node(context.Background(), &types.QueryNodeRequest{
				NodeId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetNodesCmd returns the command to query all nodes
func GetNodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Query all nodes",
		Long: fmt.Sprintf(`Query all nodes

Example:
$ %s query %s nodes --owner <owner>
`, version.AppName, types.ModuleName),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			request := &types.QueryNodesRequest{
				Pagination: pageReq,
			}

			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			if len(owner) > 0 {
				if _, err := sdk.AccAddressFromBech32(owner); err != nil {
					return err
				}
				request.Owner = owner
			}

			res, err := queryClient.Nodes(context.Background(), request)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nodes")
	cmd.Flags().String(FlagOwner, "", "The owner of nodes")
	return cmd
}

// GetDivisionCmd returns the command to query a division
func GetDivisionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "division [division-id]",
		Short: "Query the division details",
		Long: fmt.Sprintf(`Query the division details

Example:
$ %s query %s division <division-id>
`, version.AppName, types.ModuleName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Division(context.Background(),
				&types.QueryDivisionRequest{
					DivisionId: args[0],
				},
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

// GetDivisionsCmd returns the command to query all divisions
func GetDivisionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "divisions",
		Short: "Query all divisions",
		Long: fmt.Sprintf(`Query the division details

Example:
$ %s query %s divisions
`, version.AppName, types.ModuleName),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Divisions(context.Background(),
				&types.QueryDivisionsRequest{
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "divisions")
	return cmd
}

// GetSupplyCmd returns the command to query the supply of a division
func GetSupplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply [division-id]",
		Short: "Query supply of division",
		Long: fmt.Sprintf(`Query supply of division

Example:
$ %s query %s supply <division-id>
`, version.AppName, types.ModuleName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Supply(context.Background(),
				&types.QuerySupplyRequest{
					DivisionId: args[0],
				},
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

// GetSaleLevelCmd returns the command to query the sale level details
func GetSaleLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sale-level",
		Short: "Query the sale level details",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.SaleLevel(context.Background(), &types.QuerySaleLevelRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetAuthorizedMembersCmd returns the command to query the authorized members
func GetAuthorizedMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authorized-members",
		Short: "Query the authorized members",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AuthorizedMembers(context.Background(), &types.QueryAuthorizedMembersRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
