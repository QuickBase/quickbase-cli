package cmd

import (
	"github.com/spf13/cobra"
)

var relationshipCmd = &cobra.Command{
	Use:   "relationship",
	Short: "Relationship resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(relationshipCmd)
}
