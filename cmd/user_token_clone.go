package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userTokenCloneCfg *viper.Viper

var userTokenCloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone the configured user token",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		return globalCfg.Validate()
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CloneUserTokenInput{}
		qbcli.GetOptions(ctx, logger, input, userTokenCloneCfg)

		output, err := qb.CloneUserToken(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	userTokenCloneCfg, flags = cliutil.AddCommand(userTokenCmd, userTokenCloneCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CloneUserTokenInput{})
}
