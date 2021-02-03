package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userTokenDeactivateCfg *viper.Viper

var userTokenDeactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate a user token",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			qbcli.SetOptionFromArg(userTokenDeactivateCfg, args, 0, "token")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeactivateUserTokenInput{}
		qbcli.GetOptions(ctx, logger, input, userTokenDeactivateCfg)

		output, err := qb.DeactivateUserToken(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	userTokenDeactivateCfg, flags = cliutil.AddCommand(userTokenCmd, userTokenDeactivateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeactivateUserTokenInput{})
}
