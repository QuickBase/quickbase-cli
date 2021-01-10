package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/spf13/cobra"
)

// configDumpCmd represents the app command
var configDumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Print the current configuration",
	Args:  configDumpCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		config := qbclient.ConfigFileProfile{
			RealmHostname:  globalCfg.RealmHostname(),
			UserToken:      qbclient.MaskUserTokenString(globalCfg.UserToken()),
			TemporaryToken: qbclient.MaskUserTokenString(globalCfg.TemporaryToken()),
			AppID:          globalCfg.DefaultAppID(),
			TableID:        globalCfg.DefaultTableID(),
			FieldID:        globalCfg.DefaultFieldID(),
		}

		qbcli.Render(ctx, logger, cmd, globalCfg, config, nil)
	},
}

func init() {
	configCmd.AddCommand(configDumpCmd)
}

func configDumpCmdValidate(cmd *cobra.Command, args []string) error {
	return globalCfg.Validate()
}
