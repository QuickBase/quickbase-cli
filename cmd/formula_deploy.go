package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var formulaDeployCfg *viper.Viper

var formulaDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy formulas to an app",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		err = globalCfg.Validate()
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbcli.DeployFormulaInput{}
		qbcli.GetOptions(ctx, logger, input, formulaDeployCfg)

		output, err := qbcli.DeployFormula(qb, input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	formulaDeployCfg, flags = cliutil.AddCommand(formulaCmd, formulaDeployCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbcli.DeployFormulaInput{})
}
