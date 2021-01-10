package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableDeleteCfg *viper.Viper

// tableDeleteCmd represents the app get command
var tableDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a table",
	Args:  tableDeleteCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteTableInput{
			AppID:   tableDeleteCfg.GetString(qbclient.OptionAppID),
			TableID: tableDeleteCfg.GetString(qbclient.OptionTableID),
		}

		output, err := qb.DeleteTable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableDeleteCfg, flags = cliutil.AddCommand(tableCmd, tableDeleteCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionAppID, "", "", qbcli.OptionAppIDDescription)
	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
}

func tableDeleteCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default App ID if configured.
	if appID := globalCfg.DefaultAppID(); appID != "" {
		tableDeleteCfg.SetDefault(qbclient.OptionAppID, appID)
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		tableDeleteCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		tableDeleteCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	return nil
}
