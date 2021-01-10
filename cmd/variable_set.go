package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var variableSetCfg *viper.Viper

// variableSetCmd represents the app get command
var variableSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a variable",
	Args:  variableSetCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.SetVariableInput{}
		qbcli.GetOptions(ctx, logger, input, variableSetCfg)

		output, err := qb.SetVariable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	variableSetCfg, flags = cliutil.AddCommand(variableCmd, variableSetCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.SetVariableInput{})
}

func variableSetCmdValidate(cmd *cobra.Command, args []string) (err error) {
	if err = globalCfg.Validate(); err == nil {
		globalCfg.SetDefaultAppID(variableSetCfg)
		qbcli.SetOptionFromArg(variableSetCfg, args, 0, "name")
		qbcli.SetOptionFromArg(variableSetCfg, args, 1, "value")
	}
	return
}
