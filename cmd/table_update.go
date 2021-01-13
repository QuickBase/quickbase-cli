package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableUpdateCfg *viper.Viper

// tableUpdateCmd represents the app get command
var tableUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a table",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(tableUpdateCfg)
			globalCfg.SetDefaultTableID(tableUpdateCfg)
			qbcli.SetOptionFromArg(tableUpdateCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateTableInput{}
		qbcli.GetOptions(ctx, logger, input, tableUpdateCfg)

		output, err := qb.UpdateTable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableUpdateCfg, flags = cliutil.AddCommand(tableCmd, tableUpdateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.UpdateTableInput{})
}
