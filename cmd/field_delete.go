package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldDeleteCfg *viper.Viper

var fieldDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete one or many fields in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fieldDeleteCfg)
			qbcli.SetOptionFromArg(fieldDeleteCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(fieldDeleteCfg, args, 1, qbclient.OptionFieldID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteFieldsInput{}
		qbcli.GetOptions(ctx, logger, input, fieldDeleteCfg)

		output, err := qb.DeleteFields(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldDeleteCfg, flags = cliutil.AddCommand(fieldCmd, fieldDeleteCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeleteFieldsInput{})
}
