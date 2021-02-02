package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relationshipCreateCfg *viper.Viper

var relationshipCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a relationship between tables",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableIDs(relationshipCreateCfg, "child-table-id")
			qbcli.SetOptionFromArg(relationshipCreateCfg, args, 0, "child-table-id")
			qbcli.SetOptionFromArg(relationshipCreateCfg, args, 1, "parent-table-id")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateRelationshipInput{ForeignKeyField: &qbclient.CreateRelationshipInputForeignKeyField{}}
		qbcli.GetOptions(ctx, logger, input, relationshipCreateCfg)

		output, err := qb.CreateRelationship(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	relationshipCreateCfg, flags = cliutil.AddCommand(relationshipCmd, relationshipCreateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CreateRelationshipInput{ForeignKeyField: &qbclient.CreateRelationshipInputForeignKeyField{}})
}
