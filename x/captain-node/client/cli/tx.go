package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/tabilabs/tabi/x/captain-node/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	captionNodeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "nft transactions subcommands",
		Long:                       "Provides the most common nft logic for upper-level applications, compatible with Ethereum's erc721 contract",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	captionNodeTxCmd.AddCommand()

	return captionNodeTxCmd
}
