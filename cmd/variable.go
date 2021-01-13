package cmd

import (
	"github.com/spf13/cobra"
)

var variableCmd = &cobra.Command{
	Use:     "variable",
	Aliases: []string{"var"},
	Short:   "Variable resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(variableCmd)
}
