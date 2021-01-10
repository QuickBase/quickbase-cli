package cmd

import (
	"strconv"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldDeleteCfg *viper.Viper

// fieldDeleteCmd represents the app get command
var fieldDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete one or many fields in a table",
	Args:  fieldDeleteCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		fids, err := qbcli.ParseFieldList(fieldDeleteCfg.GetString(qbclient.OptionFieldID))
		cliutil.HandleError(cmd, err, "field-id option invalid")

		input := &qbclient.DeleteFieldsInput{
			TableID:  fieldDeleteCfg.GetString(qbclient.OptionTableID),
			FieldIDs: fids,
		}

		output, err := qb.DeleteFields(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldDeleteCfg, flags = cliutil.AddCommand(fieldCmd, fieldDeleteCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.String(qbclient.OptionFieldID, "", "", qbcli.OptionFieldIDDescription)
}

func fieldDeleteCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		fieldDeleteCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Set the default Field ID if configured.
	if fieldID := globalCfg.DefaultFieldID(); fieldID != 0 {
		fieldDeleteCfg.SetDefault(qbclient.OptionFieldID, strconv.Itoa(fieldID))
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		fieldDeleteCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	// Use args[1] as value for the Field ID.
	if len(args) > 1 {
		fieldDeleteCfg.SetDefault(qbclient.OptionFieldID, args[1])
	}

	return nil
}
