package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableDeleteCfg *viper.Viper

var tableDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a table",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(tableDeleteCfg)
			globalCfg.SetDefaultTableID(tableDeleteCfg)
			qbcli.SetOptionFromArg(tableDeleteCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteTableInput{}
		qbcli.GetOptions(ctx, logger, input, tableDeleteCfg)

		output, err := qb.DeleteTable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableDeleteCfg, flags = cliutil.AddCommand(tableCmd, tableDeleteCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeleteTableInput{})
}
