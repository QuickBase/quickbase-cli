package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reportListCfg *viper.Viper

var reportListCmd = &cobra.Command{
	Use:   "list",
	Short: "List reports in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultAppID(reportListCfg)
			qbcli.SetOptionFromArg(reportListCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.ListReportsInput{}
		qbcli.GetOptions(ctx, logger, input, reportListCfg)

		output, err := qb.ListReports(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	reportListCfg, flags = cliutil.AddCommand(reportCmd, reportListCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.ListReportsInput{})
}
