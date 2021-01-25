package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reportGetCfg *viper.Viper

var reportGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a report in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(reportGetCfg)
			qbcli.SetOptionFromArg(reportGetCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(reportGetCfg, args, 1, "report-id")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.GetReportInput{}
		qbcli.GetOptions(ctx, logger, input, reportGetCfg)

		output, err := qb.GetReport(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	reportGetCfg, flags = cliutil.AddCommand(reportCmd, reportGetCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.GetReportInput{})
}
