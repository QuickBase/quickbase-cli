package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appCreateCfg *viper.Viper

var appCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an app",

	Args: func(cmd *cobra.Command, args []string) error {
		return globalCfg.Validate()
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateAppInput{}
		qbcli.GetOptions(ctx, logger, input, appCreateCfg)

		output, err := qb.CreateApp(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appCreateCfg, flags = cliutil.AddCommand(appCmd, appCreateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CreateAppInput{})
}
