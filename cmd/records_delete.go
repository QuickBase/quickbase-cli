package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordsDeleteCfg *viper.Viper

var recordsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete records in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(recordsDeleteCfg)
			qbcli.SetOptionFromArg(recordsDeleteCfg, args, 0, qbclient.OptionTableID)
			recordsDeleteCfg.SetDefault("from", recordsDeleteCfg.GetString(qbclient.OptionTableID))
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteRecordsInput{}
		qbcli.GetOptions(ctx, logger, input, recordsDeleteCfg)

		output, err := qb.DeleteRecords(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	recordsDeleteCfg, flags = cliutil.AddCommand(recordsCmd, recordsDeleteCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeleteRecordsInput{})
}
