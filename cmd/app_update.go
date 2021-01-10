package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appUpdateCfg *viper.Viper

var appUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an app",

	Args: func(cmd *cobra.Command, args []string) error {
		err := globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultAppID(appUpdateCfg)
			qbcli.SetOptionFromArg(appUpdateCfg, args, 0, qbclient.OptionAppID)
		}
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateAppInput{}
		qbcli.GetOptions(ctx, logger, input, appUpdateCfg)

		output, err := qb.UpdateApp(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appUpdateCfg, flags = cliutil.AddCommand(appCmd, appUpdateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.UpdateAppInput{})
}
