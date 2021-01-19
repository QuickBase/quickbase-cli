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

var fieldOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open a field's settings page in a browser",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fieldOpenCfg)
			qbcli.SetOptionFromArg(fieldOpenCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(fieldOpenCfg, args, 1, qbclient.OptionFieldID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		opts := &FieldOpenOpts{}
		qbcli.GetOptions(ctx, logger, opts, fieldOpenCfg)

		u := fmt.Sprintf("https://%s/db/%s?a=mf&fid=%v&chain=1", globalCfg.RealmHostname(), url.PathEscape(opts.TableID), opts.FieldID)
		ctx = cliutil.ContextWithLogTag(ctx, "url", u)

		err := browser.OpenURL(u)
		logger.FatalIfError(ctx, "error opening field in browser", err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldOpenCfg, flags = cliutil.AddCommand(fieldCmd, fieldOpenCmd, qbclient.EnvPrefix)
	flags.SetOptions(&FieldOpenOpts{})
}

// FieldOpenOpts contains the options for the field open command.
type FieldOpenOpts struct {
	TableID string `cliutil:"option=table-id"`
	FieldID string `cliutil:"option=field-id"`
}
