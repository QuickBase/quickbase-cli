package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appGetCfg *viper.Viper

var appGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an app definition",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		err = globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultAppID(appGetCfg)
			qbcli.SetOptionFromArg(appGetCfg, args, 0, qbclient.OptionAppID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.GetAppInput{}
		qbcli.GetOptions(ctx, logger, input, appGetCfg)

		output, err := qb.GetApp(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appGetCfg, flags = cliutil.AddCommand(appCmd, appGetCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.GetAppInput{})
}
