package cmd

import (
	"github.com/spf13/cobra"
)

var formulaCmd = &cobra.Command{
	Use:   "formula",
	Short: "Formula resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(formulaCmd)
}
