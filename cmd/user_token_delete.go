package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var userTokenDeleteCfg *viper.Viper

var userTokenDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user token",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			qbcli.SetOptionFromArg(userTokenDeleteCfg, args, 0, "token")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteUserTokenInput{}
		qbcli.GetOptions(ctx, logger, input, userTokenDeleteCfg)

		output, err := qb.DeleteUserToken(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	userTokenDeleteCfg, flags = cliutil.AddCommand(userTokenCmd, userTokenDeleteCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeleteUserTokenInput{})
}
