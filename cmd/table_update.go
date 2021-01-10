package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableUpdateCfg *viper.Viper

// tableUpdateCmd represents the app get command
var tableUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a table",
	Args:  tableUpdateCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateTableInput{
			AppID:        tableUpdateCfg.GetString(qbclient.OptionAppID),
			TableID:      tableUpdateCfg.GetString(qbclient.OptionTableID),
			Name:         tableUpdateCfg.GetString("name"),
			Description:  tableUpdateCfg.GetString("description"),
			IconName:     tableUpdateCfg.GetString("icon-name"),
			SingularNoun: tableUpdateCfg.GetString("singular-noun"),
			PluralNoun:   tableUpdateCfg.GetString("plural-noun"),
		}

		output, err := qb.UpdateTable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableUpdateCfg, flags = cliutil.AddCommand(tableCmd, tableUpdateCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionAppID, "", "", qbcli.OptionAppIDDescription)
	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.String("name", "", "", "name for the table")
	flags.String("description", "", "", "description for the table")
	flags.String("icon-name", "", "", "icon for the table")
	flags.String("singular-noun", "", "", "singular noun for records in the table")
	flags.String("plural-noun", "", "", "plural noun for records in the table")
}

func tableUpdateCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default App ID and Table ID if configured.
	if appID := globalCfg.DefaultAppID(); appID != "" {
		tableUpdateCfg.SetDefault(qbclient.OptionAppID, appID)
	}
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		tableUpdateCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		tableUpdateCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	return nil
}
