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

var tableOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a table's page in a browser",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		err = globalCfg.Validate()
		if err == nil {
			globalCfg.SetDefaultTableID(tableOpenCfg)
			qbcli.SetOptionFromArg(tableOpenCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		opts := &TableOpenOpts{}
		qbcli.GetOptions(ctx, logger, opts, tableOpenCfg)

		u := fmt.Sprintf("https://%s/db/%s%s", globalCfg.RealmHostname(), url.PathEscape(opts.TableID), opts.Vars())
		ctx = cliutil.ContextWithLogTag(ctx, "url", u)

		err := browser.OpenURL(u)
		qbcli.HandleError(ctx, logger, "error opening table in browser", err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableOpenCfg, flags = cliutil.AddCommand(tableCmd, tableOpenCmd, qbclient.EnvPrefix)
	flags.SetOptions(&TableOpenOpts{})
}

// TableOpenOpts contains the options for the table open command.
type TableOpenOpts struct {
	TableID  string `cliutil:"option=table-id"`
	Settings string `cliutil:"option=settings"`
}

// Vars returns the query string variables based on TableOpenOpts.Settings.
func (o TableOpenOpts) Vars() string {
	if !tableOpenCfg.IsSet("settings") {
		return ""
	}

	a := "?a="
	switch o.Settings {
	case "access", "permissions":
		a += "TablePermissions"
	case "actions":
		a += "QuickBaseActionList"
	case "advanced":
		a += "KeyProps"
	case "forms":
		a += "DFormList"
	case "fields":
		a += "listfields"
	case "notifications":
		a += "EmailList"
	case "relationships":
		a += "Relationships"
	case "reports":
		a += "ReportList"
	case "webhooks":
		a += "WebhookList"
	default:
		a += "TableSettingsHome"
	}

	return a
}
