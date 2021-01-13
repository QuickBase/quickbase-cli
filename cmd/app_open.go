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

	Args: func(cmd *cobra.Command, args []string) (err error) {
		err = globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultAppID(appOpenCfg)
			qbcli.SetOptionFromArg(appOpenCfg, args, 0, qbclient.OptionAppID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		opts := &AppOpenOpts{}
		qbcli.GetOptions(ctx, logger, opts, appOpenCfg)

		u := fmt.Sprintf("https://%s/db/%s%s", globalCfg.RealmHostname(), url.PathEscape(opts.AppID), opts.Vars())
		ctx = cliutil.ContextWithLogTag(ctx, "url", u)

		err := browser.OpenURL(u)
		qbcli.HandleError(ctx, logger, "error opening app in browser", err)
	},
}

func init() {
	var flags *cliutil.Flagger
	appOpenCfg, flags = cliutil.AddCommand(appCmd, appOpenCmd, qbclient.EnvPrefix)
	flags.SetOptions(&AppOpenOpts{})
}

// AppOpenOpts contains the options for the app open command.
type AppOpenOpts struct {
	AppID    string `cliutil:"option=app-id"`
	Settings string `cliutil:"option=settings short=s"`
}

// Vars returns the query string variables based on AppOpenOpts.Settings.
func (o AppOpenOpts) Vars() string {
	if !appOpenCfg.IsSet("settings") {
		return ""
	}

	a := "?a="
	switch o.Settings {
	case "branding":
		a += "AppPropertiesBrandGuide"
	case "management", "mngmt", "mgmt", "mgt":
		a += "AppMisc"
	case "pages":
		a += "AppDBPages"
	case "properties", "props":
		a += "AppProperties"
	case "roles":
		a += "AppRolesList"
	case "tables":
		a += "AppTablesList"
	case "variables", "vars":
		a += "AppVariables"
	default:
		a += "AppSettingsHome"
	}

	return a
}
