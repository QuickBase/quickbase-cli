package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableExportCfg *viper.Viper

var tableExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export data in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(tableExportCfg)
			qbcli.SetOptionFromArg(tableExportCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		opts := &qbcli.ExportOptions{}
		qbcli.GetOptions(ctx, logger, opts, tableExportCfg)

		err := qbcli.Export(qb, opts)
		qbcli.HandleError(ctx, logger, "error exporting records", err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableExportCfg, flags = cliutil.AddCommand(tableCmd, tableExportCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbcli.ExportOptions{})
}
