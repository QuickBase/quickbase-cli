package cmd

import (
	"fmt"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var userTokenOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the My User Tokens page in a browser",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		return globalCfg.Validate()
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		u := fmt.Sprintf("https://%s/db/main?a=GenUserTokenList", globalCfg.RealmHostname())
		ctx = cliutil.ContextWithLogTag(ctx, "url", u)

		err := browser.OpenURL(u)
		qbcli.HandleError(ctx, logger, "error opening user token list in browser", err)
	},
}

func init() {
	cliutil.AddCommand(userTokenCmd, userTokenOpenCmd, qbclient.EnvPrefix)
}
