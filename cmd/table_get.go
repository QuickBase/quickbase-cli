package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableGetCfg *viper.Viper

var tableGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a table definition",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(tableGetCfg)
			globalCfg.SetDefaultTableID(tableGetCfg)
			qbcli.SetOptionFromArg(tableGetCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.GetTableInput{}
		qbcli.GetOptions(ctx, logger, input, tableGetCfg)

		output, err := qb.GetTable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableGetCfg, flags = cliutil.AddCommand(tableCmd, tableGetCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.GetTableInput{})
}
