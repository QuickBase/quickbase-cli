package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pageUpdateCfg *viper.Viper

var pageUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a page",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(pageUpdateCfg)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdatePageInput{Body: &qbclient.UpdatePageInputBody{}}
		qbcli.GetOptions(ctx, logger, input, pageUpdateCfg)

		output, err := qb.UpdatePage(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	pageUpdateCfg, flags = cliutil.AddCommand(pageCmd, pageUpdateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.UpdatePageInput{Body: &qbclient.UpdatePageInputBody{}})
}
