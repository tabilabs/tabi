package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/tabilabs/tabi/x/captains/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	captionNodeTxCmd := &cobra.Command{
		Use:                        "captains",
		Short:                      "captains transactions subcommands",
		Long:                       "Provides the most common captains logic",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	captionNodeTxCmd.AddCommand(
		NewTxCmdCreateNode(),
		NewTxCmdCommitReport(),
		NewTxCmdAddAuthorizedMembers(),
		NewTxCmdRemoveAuthorizedMembers(),
		NewTxCmdUpdateSaleLevel(),
		NewTxCmdCommitComputingPower(),
		NewTxCmdClaimComputingPower(),
		NewTxCmdDraftReport(),
	)

	return captionNodeTxCmd
}

// NewTxCmdCreateNode returns a command to mint a new Node
func NewTxCmdCreateNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-node [division-id] [receiver] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "Create a new Node and set the owner to the receiver",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
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

			msg := types.NewMsgCreateCaptainNode(sender, receiver, divisionID)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewTxCmdCommitReport returns a command to commit a report
func NewTxCmdCommitReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-report [report.json] --report-type [digest,batch,end] --from [sender]",
		Args:  cobra.ExactArgs(1),
		Short: "commit report for an epoch",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
			reportType, err := cmd.Flags().GetString(FlagReportType)
			if err != nil {
				return err
			}

			contents, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			report, err := parseReport(contents, reportType)
			if err != nil {
				return err
			}

			msg, err := types.NewMsgCommitReport(sender, parseReportType(reportType), report)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagReportType, "", "report type [digest,batch,end]")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewTxCmdAddAuthorizedMembers returns a command to add members to the authorized members list
func NewTxCmdAddAuthorizedMembers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-authorized-members [member1,member2,member3] --from [sender]",
		Args:  cobra.ExactArgs(1),
		Short: "Add members to the authorized members list",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
			members := strings.Split(args[0], ",")
			if len(members) == 0 {
				panic("members cannot be empty")
			}

			msg := types.NewAddAuthorizedMembers(sender, members)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

// NewTxCmdRemoveAuthorizedMembers returns a command to remove members from the authorized members list
func NewTxCmdRemoveAuthorizedMembers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-authorized-members [member1,member2,member3] --from [sender]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove members from the authorized members list",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()
			members := strings.Split(args[0], ",")
			if len(members) == 0 {
				panic("members cannot be empty")
			}

			msg := types.NewMsgRemoveAuthorizedMembers(sender, members)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

// NewTxCmdUpdateSaleLevel returns a command to update the sale level of a node
func NewTxCmdUpdateSaleLevel() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-sale-level [level] --from [sender]",
		Args:  cobra.ExactArgs(2),
		Short: "Update the sale level of a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := clientCtx.GetFromAddress().String()

			levelStr := strings.TrimSpace(args[1])
			if len(levelStr) == 0 {
				panic("level cannot be empty")
			}
			// convert levelStr to uint64
			level, err := strconv.ParseUint(levelStr, 10, 64)
			if err != nil {
				panic(err)
			}

			msg := types.NewMsgUpdateSaleLevel(sender, level)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

// NewTxCmdCommitComputingPower returns a command to commit computing power
func NewTxCmdCommitComputingPower() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-computing-power [path/to/update_user_experience.json] --from [sender]",
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

			sender := clientCtx.GetFromAddress().String()
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

// NewTxCmdClaimComputingPower returns a command to claim computing power
func NewTxCmdClaimComputingPower() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-computing-power [node-id] [computing-power-amount] --from [sender]",
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

			sender := clientCtx.GetFromAddress().String()

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

// NewTxCmdDraftReport returns a command to draft a report
func NewTxCmdDraftReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draft-report [digest,batch,end]",
		Args:  cobra.ExactArgs(1),
		Short: "Generate a draft report json file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return draftReport(args[0])
		},
	}
	return cmd
}
