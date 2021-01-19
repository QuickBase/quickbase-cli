package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldListCfg *viper.Viper

var fieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "List fields in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fieldListCfg)
			qbcli.SetOptionFromArg(fieldListCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.ListFieldsInput{}
		qbcli.GetOptions(ctx, logger, input, fieldListCfg)

		output, err := qb.ListFields(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldListCfg, flags = cliutil.AddCommand(fieldCmd, fieldListCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.ListFieldsInput{})
}
