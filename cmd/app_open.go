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

var appOpenCfg *viper.Viper

var appOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open an app's homepage in a browser",

	Args: func(cmd *cobra.Command, args []string) error {
		err := globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultAppID(appOpenCfg)
			qbcli.SetOptionFromArg(appOpenCfg, args, 0, qbclient.OptionAppID)
		}
		return err
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)
		appID := appOpenCfg.GetString(qbclient.OptionAppID)
		url := fmt.Sprintf("https://%s/db/%s", globalCfg.RealmHostname(), url.PathEscape(appID))
		err := browser.OpenURL(url)
		logger.FatalIfError(ctx, "error opening app in browser", err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appOpenCfg, flags = cliutil.AddCommand(appCmd, appOpenCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionAppID, "", "", qbcli.OptionAppIDDescription)
}
