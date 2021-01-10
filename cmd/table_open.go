package cmd

import (
	"fmt"
	"net/url"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableOpenCfg *viper.Viper

// tableOpenCmd represents the app get command
var tableOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a table's page in a browser",
	Args:  tableOpenCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)
		tableID := tableOpenCfg.GetString(qbclient.OptionTableID)

		var a string
		if !tableOpenCfg.GetBool("settings") {
			a = "td"
		} else {
			a = "TableSettingsHome"
		}

		u := fmt.Sprintf("https://%s/db/%s?a=%s", globalCfg.RealmHostname(), url.PathEscape(tableID), a)
		err := browser.OpenURL(u)
		logger.FatalIfError(ctx, "error opening table in browser", err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableOpenCfg, flags = cliutil.AddCommand(tableCmd, tableOpenCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.Bool("settings", "", false, "open the settings page")
}

func tableOpenCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		tableOpenCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		tableOpenCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	return nil
}
