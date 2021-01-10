package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appCopyCfg *viper.Viper

var appCopyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy an app",

	Args: func(cmd *cobra.Command, args []string) error {
		err := globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultAppID(appCopyCfg)
			qbcli.SetOptionFromArg(appCopyCfg, args, 0, qbclient.OptionAppID)
		}
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CopyAppInput{Properties: &qbclient.CopyAppInputProperties{}}
		qbcli.GetOptions(ctx, logger, input, appCopyCfg)

		output, err := qb.CopyApp(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appCopyCfg, flags = cliutil.AddCommand(appCmd, appCopyCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CopyAppInput{Properties: &qbclient.CopyAppInputProperties{}})
}
