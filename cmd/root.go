package cmd

import (
	"os"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
)

var globalCfg qbcli.GlobalConfig

var rootCmd = &cobra.Command{
	Use:   "quickbase-cli",
	Short: "A command line interface to Quick Base.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the command line tool.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg := cliutil.InitConfig(qbclient.EnvPrefix)
	globalCfg = qbcli.NewGlobalConfig(rootCmd, cfg)
}
