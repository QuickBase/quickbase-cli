package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relationshipListCfg *viper.Viper

// relationshipListCmd represents the app get command
var relationshipListCmd = &cobra.Command{
	Use:   "list",
	Short: "List a table's relationships",
	Args:  relationshipListCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)
		output, err := qb.ListRelationshipsByTableID(relationshipListCfg.GetString(qbclient.OptionTableID))
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	relationshipListCfg, flags = cliutil.AddCommand(relationshipCmd, relationshipListCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
}

func relationshipListCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		relationshipListCfg.SetDefault(qbclient.OptionTableID, tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		relationshipListCfg.SetDefault(qbclient.OptionTableID, args[0])
	}

	return nil
}
