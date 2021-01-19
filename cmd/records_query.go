package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordsQueryCfg *viper.Viper

var recordsQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query records in a table using the Quick Base query language",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(recordsQueryCfg)
			qbcli.SetOptionFromArg(recordsQueryCfg, args, 0, qbclient.OptionTableID)
			recordsQueryCfg.SetDefault("from", recordsQueryCfg.GetString(qbclient.OptionTableID))
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.QueryRecordsInput{}
		qbcli.GetOptions(ctx, logger, input, recordsQueryCfg)

		output, err := qb.QueryRecords(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	recordsQueryCfg, flags = cliutil.AddCommand(recordsCmd, recordsQueryCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.QueryRecordsInput{})
}
