package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appListCfg *viper.Viper

var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List apps",

	Args: func(cmd *cobra.Command, args []string) error {
		return globalCfg.Validate()
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.ListAppsInput{}
		qbcli.GetOptions(ctx, logger, input, appListCfg)

		output, err := qb.ListApps(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appListCfg, flags = cliutil.AddCommand(appCmd, appListCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.ListAppsInput{})
}
