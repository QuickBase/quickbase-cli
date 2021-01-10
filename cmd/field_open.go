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

var fieldOpenCfg *viper.Viper

// fieldOpenCmd represents the app get command
var fieldOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Modify a field in a browser",
	Args:  fieldOpenCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		tableID := fieldOpenCfg.GetString(qbclient.OptionTableID)
		fieldID := fieldOpenCfg.GetInt(qbclient.OptionFieldID)

		url := fmt.Sprintf("https://%s/db/%s?a=mf&fid=%v&chain=1", globalCfg.RealmHostname(), url.PathEscape(tableID), fieldID)
		err := browser.OpenURL(url)
		logger.FatalIfError(ctx, "error opening app in browser", err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldOpenCfg, flags = cliutil.AddCommand(fieldCmd, fieldOpenCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.Int(qbclient.OptionFieldID, "", 0, qbcli.OptionFieldIDDescription)
}

func fieldOpenCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		fieldOpenCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Set the default Field ID if configured.
	if fieldID := globalCfg.DefaultFieldID(); fieldID != 0 {
		fieldOpenCfg.SetDefault(qbclient.OptionFieldID, fieldID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		fieldOpenCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	// Use args[1] as value for the Field ID.
	if len(args) > 1 {
		fieldOpenCfg.SetDefault(qbclient.OptionFieldID, args[1])
	}

	return nil
}
