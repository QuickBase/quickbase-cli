package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableCreateCfg *viper.Viper

var tableCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(tableCreateCfg)
			globalCfg.SetDefaultTableID(tableCreateCfg)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateTableInput{}
		qbcli.GetOptions(ctx, logger, input, tableCreateCfg)

		output, err := qb.CreateTable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableCreateCfg, flags = cliutil.AddCommand(tableCmd, tableCreateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CreateTableInput{})
}
