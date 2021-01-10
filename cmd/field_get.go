package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldGetCfg *viper.Viper

// fieldGetCmd represents the app get command
var fieldGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a field definition",
	Args:  fieldGetCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.GetFieldInput{
			TableID: fieldGetCfg.GetString(qbclient.OptionTableID),
			FieldID: fieldGetCfg.GetInt(qbclient.OptionFieldID),
		}

		output, err := qb.GetField(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldGetCfg, flags = cliutil.AddCommand(fieldCmd, fieldGetCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.Int(qbclient.OptionFieldID, "", 0, qbcli.OptionFieldIDDescription)
}

func fieldGetCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default TableID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		fieldGetCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Set the default FieldID if configured.
	if fieldID := globalCfg.DefaultFieldID(); fieldID != 0 {
		fieldGetCfg.SetDefault(qbclient.OptionFieldID, fieldID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		fieldGetCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	// Use args[1] as value for the Field ID.
	if len(args) > 1 {
		fieldGetCfg.SetDefault(qbclient.OptionFieldID, args[1])
	}

	return nil
}
