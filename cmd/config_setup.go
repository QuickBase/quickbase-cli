package cmd

import (
	"fmt"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
)

// configSetupCmd represents the app command
var configSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run initial setup of a configuration file",
	Args:  configSetupCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {

		// Bail if there is an existing configuration file.
		filepath := qbclient.Filepath(globalCfg.ConfigDir(), qbclient.ConfigFilename)
		if qbclient.FileExists(filepath) {
			fmt.Printf("Configuration file already created at %s.\nPlease edit the file directly.\n", filepath)
			return
		}

		hostname, err := qbcli.Prompt("Realm Hostname: ", qbclient.ValidateHostname)
		cliutil.HandleError(cmd, err, "error reading realm hostname")

		usertoken, err := qbcli.Prompt("User Token: ", qbclient.ValidateNotEmptyFn("user token"))
		cliutil.HandleError(cmd, err, "error reading user token")

		appID, err := qbcli.Prompt("App ID (optional): ", qbclient.NoValidation)
		cliutil.HandleError(cmd, err, "error reading app id")

		cf := make(map[string]*qbclient.ConfigFileProfile, 1)
		cf["default"] = &qbclient.ConfigFileProfile{
			RealmHostname: hostname,
			UserToken:     usertoken,
			AppID:         appID,
		}

		err = qbclient.WriteConfigFile(globalCfg.ConfigDir(), cf)
		cliutil.HandleError(cmd, err, "error writing config file")

		fmt.Printf("\nConfig file written to %s.\n", filepath)
	},
}

func init() {
	configCmd.AddCommand(configSetupCmd)
}

func configSetupCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
