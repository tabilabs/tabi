package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/tabilabs/tabi/x/captains/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	captionNodeTxCmd := &cobra.Command{
		Use:                        "captain-node",
		Short:                      "captain-node transactions subcommands",
		Long:                       "Provides the most common captain-node logic",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	captionNodeTxCmd.AddCommand(
		NewMintNodeTxCmd(),
		NewUpdateUserExperienceCmd(),
		NewCommitReportCmd(),
		NewWithdrawExperienceCmd(),
	)

	return captionNodeTxCmd
}

// NewMintNodeTxCmd returns a command to mint a new Node
func NewMintNodeTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [division-id] [receiver] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "Mint a new Node and set the owner to the receiver",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s mint <division-id> <receiver> --from <sender> --chain-id <chain-id>`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()
			divisionID := strings.TrimSpace(args[0])
			if len(divisionID) == 0 {
				panic("division-id cannot be empty")
			}

			receiver := strings.TrimSpace(args[1])
			if len(receiver) > 0 {
				if _, err = sdk.AccAddressFromBech32(receiver); err != nil {
					return err
				}
			} else {
				panic("receiver cannot be empty")
			}

			msg := types.NewMsgCreateCaptainNode(divisionID, receiver, sender)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCommitReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-report [path/to/update_power_on_period.json] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "update power on period for multiple nodes",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s commit-report ./update_power_on_period.json --from <sender> --chain-id <chain-id>
Where update_power_on_period.json contains:

[
    {
      "power_on_period": "6",
      "node_id": "0x00000000000000001",
    },
	.....
],

`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()

			messages := args[0]
			if !json.Valid([]byte(messages)) {
				messagesContent, err := os.ReadFile(messages)
				if err != nil {
					return fmt.Errorf("invalid options: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(messagesContent) {
					return fmt.Errorf("invalid options: .json file content is invalid JSON")
				}

				messages = string(messagesContent)
			}

			// FIXME: we need report!
			//var captainNodePowerOnPeriods []*types.CaptainNodePowerOnPeriod
			//if err := json.Unmarshal([]byte(messages), &captainNodePowerOnPeriods); err != nil {
			//	return fmt.Errorf("failed to unmarshal JSON: %w", err)
			//}

			msg := types.NewMsgCommitReport(sender, "")
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewUpdateUserExperienceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-computing-power [path/to/update_user_experience.json] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "Mint a new Node and set the owner to the receiver",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s reward-computing-power <division-id> <receiver> --from <sender> --chain-id <chain-id>
Where update_power_on_period.json contains:
{
  // array of proto-JSON-encoded sdk.Msgs
[
    {
      "experience": "6",
      "receiver": "tabixxxxxxx",
    },
	.....
],


`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()
			messages := args[0]
			if !json.Valid([]byte(messages)) {
				messagesContent, err := os.ReadFile(messages)
				if err != nil {
					return fmt.Errorf("invalid options: neither JSON input nor path to .json file were provided")
				}

				if !json.Valid(messagesContent) {
					return fmt.Errorf("invalid options: .json file content is invalid JSON")
				}

				messages = string(messagesContent)
			}

			var extractableComputingPowers []types.ClaimableComputingPower
			if err := json.Unmarshal([]byte(messages), &extractableComputingPowers); err != nil {
				return fmt.Errorf("failed to unmarshal JSON: %w", err)
			}

			msg := types.NewMsgCommitComputingPower(extractableComputingPowers, sender)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewWithdrawExperienceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-computing-power [node-id] [computing-power-amount] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "Withdraw experience to a node",
		Long: strings.TrimSpace(fmt.Sprintf(`
			$ %s tx %s withdraw-experience <node-id> <experience-amount> --from <sender> --chain-id <chain-id>`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var sender = clientCtx.GetFromAddress().String()

			nodeID := strings.TrimSpace(args[0])
			if len(nodeID) == 0 {
				panic("node-id cannot be empty")
			}

			experienceAmountStr := strings.TrimSpace(args[1])
			if len(experienceAmountStr) == 0 {
				panic("experience-amount cannot be empty")
			}
			// convert experienceAmountStr to uint64
			experienceAmount, err := strconv.ParseUint(experienceAmountStr, 10, 64)
			if err != nil {
				panic(err)
			}

			msg := types.NewMsgWithdrawComputingPower(nodeID, experienceAmount, sender)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
