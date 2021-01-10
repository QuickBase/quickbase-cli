package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appEventsCfg *viper.Viper

var appEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List app events",

	Args: func(cmd *cobra.Command, args []string) error {
		err := globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultAppID(appEventsCfg)
			qbcli.SetOptionFromArg(appEventsCfg, args, 0, qbclient.OptionAppID)
		}
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.ListAppEventsInput{}
		qbcli.GetOptions(ctx, logger, input, appEventsCfg)

		output, err := qb.ListAppEvents(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appEventsCfg, flags = cliutil.AddCommand(appCmd, appEventsCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.ListAppEventsInput{})
}
