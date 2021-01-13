package cmd

import (
	"fmt"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(qbclient.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
