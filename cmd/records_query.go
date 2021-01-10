package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordsQueryCfg *viper.Viper

// recordsQueryCmd represents the record query command
var recordsQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query records in a table using the Quick Base query language",
	Args:  recordsQueryCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		// Parse the list of fields IDs.
		fids, err := qbcli.ParseFieldList(recordsQueryCfg.GetString("select"))
		cliutil.HandleError(cmd, err, "select option invalid")

		// Parse the "where" clause.
		where := qbcli.ParseQuery(recordsQueryCfg.GetString("where"))
		ctx = cliutil.ContextWithLogTag(ctx, "query", where)
		logger.Debug(ctx, "query formatted")

		input := &qbclient.QueryRecordsInput{
			Select: fids,
			From:   recordsQueryCfg.GetString("from"),
			Where:  where,
		}

		// Add sortBy clause if passed.
		sortBy := recordsQueryCfg.GetString("sort-by")
		if sortBy != "" {
			input.SortBy, err = qbcli.ParseSortBy(sortBy)
			cliutil.HandleError(cmd, err, "sort-by option invalid")
		}

		// Add groupBy clause if passed.
		groupBy := recordsQueryCfg.GetString("group-by")
		if groupBy != "" {
			input.GroupBy, err = qbcli.ParseGroupBy(groupBy)
			cliutil.HandleError(cmd, err, "group-by option invalid")
		}

		output, err := qb.QueryRecords(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	recordsQueryCfg, flags = cliutil.AddCommand(recordsCmd, recordsQueryCmd, qbclient.EnvPrefix)

	flags.String("select", "", "", "comma-separated list of fields returned in the response (required)")
	flags.String("from", "", "", qbcli.OptionTableIDDescription)
	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription+" (alias for \"from\")")
	flags.String("where", "", "", "filter, using the Quick Base query language, which determines the records to return (required)")
	flags.String("group-by", "", "", "fields to group the records by, e.g., \"6 same-value\"")
	flags.String("sort-by", "", "", "fields to sort the records by, e.g., \"7 DESC,8 ASC\"")
}

func recordsQueryCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		recordsQueryCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		recordsQueryCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	// This is how we handle option alises. I wish Viper's aliases worked with Cobra!
	recordsQueryCfg.SetDefault("from", recordsQueryCfg.GetString(qbclient.OptionTableID))

	return nil
}
