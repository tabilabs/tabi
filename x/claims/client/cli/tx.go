package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/spf13/cobra"
	"github.com/tabilabs/tabi/x/claims/types"
)

// NewTxCmd returns a root CLI command handler for all x/distribution transaction commands.
func NewTxCmd() *cobra.Command {
	claimsTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Distribution transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	claimsTxCmd.AddCommand(
		NewWithdrawRewardsCmd(),
	)

	return claimsTxCmd
}

// NewWithdrawRewardsCmd returns a CLI command handler for creating a MsgWithdrawDelegatorReward transaction.
func NewWithdrawRewardsCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "withdraw-rewards [node_id]",
		Short: "Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Withdraw rewards from a given delegation address,
and optionally withdraw validator commission if the delegation address given is a validator operator.

Example:
$ %s tx distribution withdraw-rewards xxxxxxx --from mykey
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			ownerAddr := clientCtx.GetFromAddress()

			msgs := []sdk.Msg{types.NewMsgWithdrawNodeReward(ownerAddr, args[0])}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
