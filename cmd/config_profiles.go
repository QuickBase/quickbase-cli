package cmd

import (
	"github.com/spf13/cobra"
)

// configProfilesCfg represents the app command
var configProfilesCfg = &cobra.Command{
	Use:   "profiles",
	Short: "List profiles in the config file",
	Args:  configProfilesCfgValidate,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	configCmd.AddCommand(configProfilesCfg)
}

func configProfilesCfgValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
