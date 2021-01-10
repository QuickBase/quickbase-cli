package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relationshipUpdateCfg *viper.Viper

// relationshipUpdateCmd represents the app get command
var relationshipUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a relationship",
	Args:  relationshipUpdateCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateRelationshipInput{
			ChildTableID:   relationshipUpdateCfg.GetString("child-table-id"),
			RelationshipID: relationshipUpdateCfg.GetInt(qbclient.OptionRelationshipID),
		}

		// Parse the list of fields IDs.
		fids, err := qbcli.ParseFieldList(relationshipUpdateCfg.GetString("lookup-fields"))
		cliutil.HandleError(cmd, err, "lookup-fields option invalid")
		if len(fids) > 0 {
			input.LookupFieldIDs = fids
		}

		output, err := qb.UpdateRelationship(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	relationshipUpdateCfg, flags = cliutil.AddCommand(relationshipCmd, relationshipUpdateCmd, qbclient.EnvPrefix)

	flags.String("child-table-id", "", "", "unique identifier (dbid) of the child table (required)")
	flags.Int(qbclient.OptionRelationshipID, "", 0, "unique identifier of the relationship (required)")
	flags.String("lookup-fields", "", "", "ids of lookup fields in the parent table")
}

func relationshipUpdateCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if tableID := globalCfg.DefaultTableID(); tableID != "" {
		relationshipUpdateCfg.SetDefault("child-table-id", tableID)
	}

	// Use args[0] as value for the Table ID.
	if len(args) > 0 {
		relationshipUpdateCfg.SetDefault("child-table-id", args[0])
	}

	// Use args[1] as value for the Relationship ID.
	if len(args) > 1 {
		relationshipUpdateCfg.SetDefault(qbclient.OptionRelationshipID, args[1])
	}

	return nil
}
