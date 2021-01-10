package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldUpdateCfg *viper.Viper

// fieldUpdateCmd represents the app get command
var fieldUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a field in a table",
	Args:  fieldUpdateCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateFieldInput{
			TableID: fieldUpdateCfg.GetString(qbclient.OptionTableID),
			Field:   qbcli.NewFieldFromOptions(fieldUpdateCfg),
			Properties: &qbclient.UpdateFieldInputProperties{
				FieldProperties: qbcli.NewPropertiesFromOptions(fieldUpdateCfg),
			},
		}

		if !fieldUpdateCfg.GetBool("dry-run") {
			output, err := qb.UpdateField(input)
			qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
		} else {
			cliutil.PrintJSON(input)
		}
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldUpdateCfg, flags = cliutil.AddCommand(fieldCmd, fieldUpdateCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.Int(qbclient.OptionFieldID, "", 0, qbcli.OptionFieldIDDescription)
	flags.Bool("dry-run", "", false, "print input but don't sent it")

	qbcli.SetFieldOptions(flags)
}

func fieldUpdateCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default TableID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		fieldUpdateCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Set the default FieldID if configured.
	if fieldID := globalCfg.DefaultFieldID(); fieldID != 0 {
		fieldUpdateCfg.SetDefault(qbclient.OptionFieldID, fieldID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		fieldUpdateCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	// Use args[1] as value for the Field ID.
	if len(args) > 1 {
		fieldUpdateCfg.SetDefault(qbclient.OptionFieldID, args[1])
	}

	return nil
}
