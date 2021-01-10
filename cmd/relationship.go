package cmd

import (
	"github.com/spf13/cobra"
)

// relationshipCmd represents the app command
var relationshipCmd = &cobra.Command{
	Use:   "relationship",
	Short: "relationship resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(relationshipCmd)
}
