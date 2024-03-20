package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/spf13/cobra"

	"github.com/tabilabs/tabi/x/captain-node/types"
)

// GetQueryCmd returns the cli query commands for the mint module.
func GetQueryCmd() *cobra.Command {
	captionNodeQueryCmd := &cobra.Command{
		Use:                        "captain-node",
		Short:                      "Querying commands for the cliams module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	captionNodeQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetOwnerCmd(),
		GetSupplyCmd(),
		GetDivisionCmd(),
		GetDivisionsCmd(),
		GetNodeCmd(),
		GetNodesCmd(),
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

func GetOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner [node-id]",
		Short: "Query the node of owner",
		Long: fmt.Sprintf(` Query the node of owner

Example:
$ %s query %s owner <node-id>
`, version.AppName, types.ModuleName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			nodeId := args[0]

			res, err := queryClient.Owner(context.Background(), &types.QueryOwnerRequest{Id: nodeId})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetSupplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply [division-id]",
		Short: "Query supply of division",
		Long: fmt.Sprintf(`Query supply of division

Example:
$ %s query %s supply <division-id>
`, version.AppName, types.ModuleName),
		Args: cobra.NoArgs,
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

func GetDivisionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "division [division-id]",
		Short: "Query the division details",
		Long: fmt.Sprintf(`Query the division details

Example:
$ %s query %s division <division-id>
`, version.AppName, types.ModuleName),
		Args: cobra.NoArgs,
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
	return cmd
}

func GetNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node [node-id]",
		Short: "Query the node details",
		Long: fmt.Sprintf(`Query the division details

Example:
$ %s query %s node <node-id>
`, version.AppName, types.ModuleName),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Node(context.Background(), &types.QueryNodeRequest{
				Id: args[0],
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

func GetNodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Query all nodes",
		Long: fmt.Sprintf(`Query all nodes

Example:
$ %s query %s node <node-id>
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

			res, err := queryClient.Nodes(context.Background(),
				&types.QueryNodesRequest{
					Pagination: pageReq,
					Owner:      args[0],
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
