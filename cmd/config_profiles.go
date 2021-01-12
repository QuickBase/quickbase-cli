package cmd

import (
	"fmt"
	"sort"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/spf13/cobra"
)

var configProfilesCfg = &cobra.Command{
	Use:   "profiles",
	Short: "List profiles in the config file",

	Args: func(cmd *cobra.Command, args []string) error {
		return globalCfg.ReadInConfig()
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, _ := qbcli.NewLogger(cmd, globalCfg)

		cfg, err := qbclient.ReadConfigFile(globalCfg.ConfigDir())
		qbcli.HandleError(ctx, logger, "error reading config file", err)

		i := 0
		profiles := make([]string, len(cfg))
		for profile := range cfg {
			profiles[i] = profile
			i++
		}

		sort.Strings(profiles)
		for _, profile := range profiles {
			fmt.Println(profile)
		}
	},
}

func init() {
	configCmd.AddCommand(configProfilesCfg)
}
