package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableGetCfg *viper.Viper

// tableGetCmd represents the app get command
var tableGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a table definition",
	Args:  tableGetCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)
		output, err := qb.GetTableByID(tableGetCfg.GetString(qbclient.OptionTableID))
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableGetCfg, flags = cliutil.AddCommand(tableCmd, tableGetCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
}

func tableGetCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		tableGetCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		tableGetCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	return nil
}
