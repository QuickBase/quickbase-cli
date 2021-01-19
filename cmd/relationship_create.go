package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relationshipCreateCfg *viper.Viper

// relationshipCreateCmd represents the app get command
var relationshipCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a relationship between tables",
	Args:  relationshipCreateCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateRelationshipInput{
			ChildTableID:  relationshipCreateCfg.GetString("child-table-id"),
			ParentTableID: relationshipCreateCfg.GetString("parent-table-id"),
		}

		if label := relationshipCreateCfg.GetString("foreign-key-label"); label == "" {
			input.ForeignKeyField = &qbclient.CreateRelationshipInputForeignKeyField{Label: label}
		}

		// Parse the list of fields IDs.
		fids, err := cliutil.ParseIntSlice(relationshipCreateCfg.GetString("lookup-fields"))
		cliutil.HandleError(cmd, err, "lookup-fields option invalid")
		if len(fids) > 0 {
			input.LookupFieldIDs = fids
		}

		output, err := qb.CreateRelationship(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	relationshipCreateCfg, flags = cliutil.AddCommand(relationshipCmd, relationshipCreateCmd, qbclient.EnvPrefix)

	flags.String("child-table-id", "", "", "unique identifier (dbid) of the child table (required)")
	flags.String("parent-table-id", "", "", "unique identifier (dbid) of the parent table (required)")
	flags.String("foreign-key-label", "", "", "label for the foreign key field")
	flags.String("lookup-fields", "", "", "ids of lookup fields in the parent table")
}

func relationshipCreateCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default Table ID if configured.
	if childTableID := globalCfg.DefaultTableID(); childTableID != "" {
		relationshipCreateCfg.SetDefault("child-table-id", childTableID)
	}

	// Use args[0] as value for the Child Table ID.
	if len(args) > 0 {
		relationshipCreateCfg.SetDefault("child-table-id", args[0])
	}

	// Use args[1] as value for the Parent Table ID.
	if len(args) > 1 {
		relationshipCreateCfg.SetDefault("parent-table-id", args[1])
	}

	return nil
}
