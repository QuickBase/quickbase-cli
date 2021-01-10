package cmd

import (
	"github.com/spf13/cobra"
)

// recordsCmd represents the app command
var recordsCmd = &cobra.Command{
	Use:   "records",
	Short: "Records resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(recordsCmd)
}
