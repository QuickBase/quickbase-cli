package cmd

import (
	"errors"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
)

var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run initial setup of a configuration file",

	Args: func(cmd *cobra.Command, args []string) error {
		return globalCfg.ReadInConfig()
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		filepath := qbclient.Filepath(globalCfg.ConfigDir(), qbclient.ConfigFilename)
		ctx = cliutil.ContextWithLogTag(ctx, "file", filepath)

		if qbclient.FileExists(filepath) {
			err := errors.New("config file already created")
			qbcli.HandleError(ctx, logger, "please edit the config file directly", err)
		}

		hostname, err := qbcli.Prompt("Realm Hostname: ", qbclient.ValidateHostname)
		qbcli.HandleError(ctx, logger, "error reading realm hostname", err)

		usertoken, err := qbcli.Prompt("User Token: ", qbclient.ValidateNotEmptyFn("user token"))
		qbcli.HandleError(ctx, logger, "error reading user token", err)

		appID, err := qbcli.Prompt("App ID (optional): ", qbclient.NoValidation)
		qbcli.HandleError(ctx, logger, "error reading app id", err)

		cf := make(map[string]*qbclient.ConfigFileProfile, 1)
		cf["default"] = &qbclient.ConfigFileProfile{
			RealmHostname: hostname,
			UserToken:     usertoken,
			AppID:         appID,
		}

		err = qbclient.WriteConfigFile(globalCfg.ConfigDir(), cf)
		qbcli.HandleError(ctx, logger, "error writing config file", err)
		logger.Notice(ctx, "setup complete")
	},
}

func init() {
	configCmd.AddCommand(configSetupCmd)
}
