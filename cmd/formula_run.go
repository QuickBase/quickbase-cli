package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var formulaRunCfg *viper.Viper

var formulaRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a formula",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		err = globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultTableID(formulaRunCfg)
			qbcli.SetOptionFromArg(formulaRunCfg, args, 0, qbclient.OptionTableID)
			formulaRunCfg.SetDefault("from", formulaRunCfg.GetString(qbclient.OptionTableID))
			qbcli.SetOptionFromArg(formulaRunCfg, args, 1, "record-id")
			qbcli.SetOptionFromArg(formulaRunCfg, args, 2, "formula")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.RunFormulaInput{}
		qbcli.GetOptions(ctx, logger, input, formulaRunCfg)

		output, err := qb.RunFormula(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	formulaRunCfg, flags = cliutil.AddCommand(formulaCmd, formulaRunCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.RunFormulaInput{})
}
