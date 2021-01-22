package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pageGetCfg *viper.Viper

var pageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a page",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(pageGetCfg)
			qbcli.SetOptionFromArg(pageGetCfg, args, 0, "page-id")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.GetPageInput{}
		qbcli.GetOptions(ctx, logger, input, pageGetCfg)

		output, err := qb.GetPage(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	pageGetCfg, flags = cliutil.AddCommand(pageCmd, pageGetCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.GetPageInput{})
}
