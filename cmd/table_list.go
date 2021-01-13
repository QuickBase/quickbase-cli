package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableListCfg *viper.Viper

var tableListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tables in an app",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(tableListCfg)
			qbcli.SetOptionFromArg(tableListCfg, args, 0, qbclient.OptionAppID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.ListTablesInput{}
		qbcli.GetOptions(ctx, logger, input, tableListCfg)

		output, err := qb.ListTables(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableListCfg, flags = cliutil.AddCommand(tableCmd, tableListCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.ListTablesInput{})
}
