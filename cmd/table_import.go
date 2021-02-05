package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableImportCfg *viper.Viper

var tableImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data into a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(tableImportCfg)
			qbcli.SetOptionFromArg(tableImportCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		opts := &qbcli.ImportOptions{}
		qbcli.GetOptions(ctx, logger, opts, tableImportCfg)

		output, err := qbcli.Import(qb, opts)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableImportCfg, flags = cliutil.AddCommand(tableCmd, tableImportCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbcli.ImportOptions{})
}
