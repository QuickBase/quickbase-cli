package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reportRunCfg *viper.Viper

var reportRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a report in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(reportRunCfg)
			qbcli.SetOptionFromArg(reportRunCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(reportRunCfg, args, 1, "report-id")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.RunReportInput{}
		qbcli.GetOptions(ctx, logger, input, reportRunCfg)

		output, err := qb.RunReport(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	reportRunCfg, flags = cliutil.AddCommand(reportCmd, reportRunCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.RunReportInput{})
}
