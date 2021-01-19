package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldUpdateCfg *viper.Viper

var fieldUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a field in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fieldUpdateCfg)
			qbcli.SetOptionFromArg(fieldUpdateCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(fieldUpdateCfg, args, 1, qbclient.OptionFieldID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateFieldInput{Properties: &qbclient.UpdateFieldInputProperties{}}
		qbcli.GetOptions(ctx, logger, input, fieldUpdateCfg)

		output, err := qb.UpdateField(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldUpdateCfg, flags = cliutil.AddCommand(fieldCmd, fieldUpdateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.UpdateFieldInput{Properties: &qbclient.UpdateFieldInputProperties{}})
}
