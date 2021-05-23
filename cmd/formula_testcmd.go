package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var formulaTestCfg *viper.Viper

var formulaTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test formulas",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		err = globalCfg.Validate()
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbcli.TestFormulaInput{}
		qbcli.GetOptions(ctx, logger, input, formulaTestCfg)

		output, err := qbcli.TestFormula(qb, input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)

		if len(output.Failed) > 0 {
			err = qbcli.TestsFailedError("%v tests failed", len(output.Failed))
			qbcli.HandleError(ctx, logger, "tests failed", err)
		}
	},
}

func init() {
	var flags *cliutil.Flagger
	formulaTestCfg, flags = cliutil.AddCommand(formulaCmd, formulaTestCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbcli.TestFormulaInput{})
}
