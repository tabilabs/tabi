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
		Short:                      "Claims transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	claimsTxCmd.AddCommand(
		NewClaimsCmd(),
	)

	return claimsTxCmd
}

// NewClaimsCmd returns a CLI command handler for creating a MsgWithdrawDelegatorReward transaction.
func NewClaimsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claims [receiver]",
		Short: "claims rewards",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Claims rewards. 
Example:
$ %s tx claims claims xxxxxxx --from mykey
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
			sender := clientCtx.GetFromAddress()
			receiver, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			msgs := []sdk.Msg{types.NewMsgClaims(sender, receiver)}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msgs...)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
