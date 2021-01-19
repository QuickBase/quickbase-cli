package cmd

import (
	"github.com/spf13/cobra"
)

var fieldCmd = &cobra.Command{
	Use:   "field",
	Short: "Fields resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(fieldCmd)
}
