package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldGetCfg *viper.Viper

var fieldGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a field definition",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fieldGetCfg)
			qbcli.SetOptionFromArg(fieldGetCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(fieldGetCfg, args, 1, qbclient.OptionFieldID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.GetFieldInput{}
		qbcli.GetOptions(ctx, logger, input, fieldGetCfg)

		output, err := qb.GetField(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldGetCfg, flags = cliutil.AddCommand(fieldCmd, fieldGetCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.GetFieldInput{})
}
