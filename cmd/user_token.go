package cmd

import (
	"github.com/spf13/cobra"
)

var userTokenCmd = &cobra.Command{
	Use:     "user-token",
	Aliases: []string{"token"},
	Short:   "User token resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(userTokenCmd)
}
