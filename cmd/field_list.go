package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldListCfg *viper.Viper

// fieldListCmd represents the app get command
var fieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "List fields in a table",
	Args:  fieldListCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.ListFieldsInput{
			TableID:                 fieldListCfg.GetString(qbclient.OptionTableID),
			IncludeFieldPermissions: fieldListCfg.GetBool("include-field-permissions"),
		}

		output, err := qb.ListFields(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldListCfg, flags = cliutil.AddCommand(fieldCmd, fieldListCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.Bool("include-field-permissions", "", false, "return custom permissions for the fields")
}

func fieldListCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		fieldListCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		fieldListCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	return nil
}
