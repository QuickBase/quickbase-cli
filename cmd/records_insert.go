package cmd

import (
	"fmt"
	"strconv"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recordsInsertCfg *viper.Viper

// recordsInsertCmd represents the app get command
var recordsInsertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Insert and/or update records in a table",
	Args:  recordsInsertCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		tableID := recordsInsertCfg.GetString("to")
		tmap, err := qbcli.GetFieldTypeMap(qb, tableID)
		logger.FatalIfError(ctx, "error getting field types", err)

		record := &qbclient.Record{}
		data := cliutil.ParseKeyValue(recordsInsertCfg.GetString("data"))
		for fidstr, val := range data {

			// TODO check for field alias
			fid, err := strconv.Atoi(fidstr)
			cliutil.HandleError(cmd, err, fmt.Sprintf("invalid fid (%s)", fidstr))

			ftype, ok := tmap[fid]
			if !ok {
				cliutil.HandleError(cmd, fmt.Errorf("field not in table (fid %s)", fidstr), "")
			}

			v, err := qbclient.NewValueFromString(val, ftype)
			cliutil.HandleError(cmd, err, fmt.Sprintf("invalid value (fid %s)", fidstr))
			record.SetValue(fid, v)
		}

		fids, err := qbcli.ParseFieldList(recordsInsertCfg.GetString("fields-to-return"))
		cliutil.HandleError(cmd, err, "invalid option (fields-to-return)")

		input := &qbclient.InsertRecordsInput{
			To:             recordsInsertCfg.GetString("to"),
			MergeFieldID:   recordsInsertCfg.GetInt("merge-field-id"),
			FieldsToReturn: fids,
		}
		input.SetRecords([]*qbclient.Record{record})

		output, err := qb.InsertRecords(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	recordsInsertCfg, flags = cliutil.AddCommand(recordsCmd, recordsInsertCmd, qbclient.EnvPrefix)

	flags.String("to", "", "", qbcli.OptionTableIDDescription)
	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription+" (alias for \"to\")")
	flags.String("data", "", "", "record data to create")
	flags.Int("merge-field-id", "", 0, "merge field id")
	flags.String("fields-to-return", "", "", "comma-separated list of fields returned in the response for any updated or added records")
}

func recordsInsertCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		recordsInsertCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		recordsInsertCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	// This is how we handle option alises. I wish Viper's aliases worked with Cobra!
	recordsInsertCfg.SetDefault("to", recordsInsertCfg.GetString(qbclient.OptionTableID))

	return nil
}
