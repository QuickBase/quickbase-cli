package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appDeleteCfg *viper.Viper

var appDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an app",

	Args: func(cmd *cobra.Command, args []string) error {
		err := globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultAppID(appDeleteCfg)
			qbcli.SetOptionFromArg(appDeleteCfg, args, 0, qbclient.OptionAppID)
		}
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteAppInput{}
		qbcli.GetOptions(ctx, logger, input, appDeleteCfg)

		output, err := qb.DeleteApp(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appDeleteCfg, flags = cliutil.AddCommand(appCmd, appDeleteCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeleteAppInput{})
}
