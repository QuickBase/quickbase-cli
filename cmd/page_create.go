package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pageCreateCfg *viper.Viper

var pageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a page",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(pageCreateCfg)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreatePageInput{Body: &qbclient.CreatePageInputBody{}}
		qbcli.GetOptions(ctx, logger, input, pageCreateCfg)

		output, err := qb.CreatePage(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	pageCreateCfg, flags = cliutil.AddCommand(pageCmd, pageCreateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CreatePageInput{Body: &qbclient.CreatePageInputBody{}})
}
