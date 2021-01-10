package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relationshipDeleteCfg *viper.Viper

// relationshipDeleteCmd represents the app get command
var relationshipDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a relationship",
	Args:  relationshipDeleteCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteRelationshipInput{
			ChildTableID:   relationshipDeleteCfg.GetString("child-table-id"),
			RelationshipID: relationshipDeleteCfg.GetInt(qbclient.OptionRelationshipID),
		}

		output, err := qb.DeleteRelationship(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	relationshipDeleteCfg, flags = cliutil.AddCommand(relationshipCmd, relationshipDeleteCmd, qbclient.EnvPrefix)

	flags.String("child-table-id", "", "", "unique identifier (dbid) of the child table (required)")
	flags.Int(qbclient.OptionRelationshipID, "", 0, "unique identifier of the relationship (required)")
}

func relationshipDeleteCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		relationshipDeleteCfg.SetDefault("child-table-id", tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		relationshipDeleteCfg.SetDefault("child-table-id", args[0])
	}

	// Use args[1] as value for the Relationship ID.
	if len(args) > 1 {
		relationshipDeleteCfg.SetDefault(qbclient.OptionRelationshipID, args[1])
	}

	return nil
}
