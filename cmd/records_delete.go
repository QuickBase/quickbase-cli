package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordsDeleteCfg *viper.Viper

// recordsDeleteCmd represents the app get command
var recordsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete records in a table",
	Args:  recordsDeleteCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		where := qbcli.ParseQuery(recordsDeleteCfg.GetString("where"))
		ctx = cliutil.ContextWithLogTag(ctx, "query", where)
		logger.Debug(ctx, "query formatted")

		input := &qbclient.DeleteRecordsInput{
			From:  recordsDeleteCfg.GetString("from"),
			Where: where,
		}

		output, err := qb.DeleteRecords(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	recordsDeleteCfg, flags = cliutil.AddCommand(recordsCmd, recordsDeleteCmd, qbclient.EnvPrefix)

	flags.String("from", "", "", qbcli.OptionTableIDDescription)
	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription+" (alias for \"from\")")
	flags.String("where", "", "", "filter, using the Quick Base query language, which determines the records to delete (required)")
}

func recordsDeleteCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		recordsDeleteCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		recordsDeleteCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	// This is how we handle option alises. I wish Viper's aliases worked with Cobra!
	recordsDeleteCfg.SetDefault("from", recordsDeleteCfg.GetString(qbclient.OptionTableID))

	return nil
}
