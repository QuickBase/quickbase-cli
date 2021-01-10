package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the app command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
