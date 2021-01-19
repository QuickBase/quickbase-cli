package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var variableGetCfg *viper.Viper

// variableGetCmd represents the app get command
var variableGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a variable",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(variableGetCfg)
			qbcli.SetOptionFromArg(variableGetCfg, args, 0, "variable-name")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.GetVariableInput{}
		qbcli.GetOptions(ctx, logger, input, variableGetCfg)

		output, err := qb.GetVariable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	variableGetCfg, flags = cliutil.AddCommand(variableCmd, variableGetCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.GetVariableInput{})
}
