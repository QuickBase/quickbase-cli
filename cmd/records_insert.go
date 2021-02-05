package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordsInsertCfg *viper.Viper

var recordsInsertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert and/or update records in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(recordsInsertCfg)
			qbcli.SetOptionFromArg(recordsInsertCfg, args, 0, qbclient.OptionTableID)
			recordsInsertCfg.SetDefault("to", recordsInsertCfg.GetString(qbclient.OptionTableID))
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		err := qbcli.CacheTableSchema(qb, recordsInsertCfg.GetString("to"))
		qbcli.HandleError(ctx, logger, "error setting field type map", err)

		input := &qbclient.InsertRecordsInput{}
		qbcli.GetOptions(ctx, logger, input, recordsInsertCfg)

		output, err := qb.InsertRecords(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	recordsInsertCfg, flags = cliutil.AddCommand(recordsCmd, recordsInsertCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.InsertRecordsInput{})
}
