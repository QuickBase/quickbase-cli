package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableListCfg *viper.Viper

// tableListCmd represents the app get command
var tableListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tables in an app",
	Args:  tableListCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)
		output, err := qb.ListTablesByAppID(tableListCfg.GetString(qbclient.OptionAppID))
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableListCfg, flags = cliutil.AddCommand(tableCmd, tableListCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionAppID, "", "", qbcli.OptionAppIDDescription)
}

func tableListCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default App ID if configured.
	if appID := globalCfg.DefaultAppID(); appID != "" {
		tableListCfg.SetDefault(qbclient.OptionAppID, appID)
	}

	return nil
}
