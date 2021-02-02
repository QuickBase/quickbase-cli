package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relationshipDeleteCfg *viper.Viper

var relationshipDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a relationship",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableIDs(relationshipDeleteCfg, "child-table-id")
			qbcli.SetOptionFromArg(relationshipDeleteCfg, args, 0, "child-table-id")
			qbcli.SetOptionFromArg(relationshipDeleteCfg, args, 1, qbclient.OptionRelationshipID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteRelationshipInput{}
		qbcli.GetOptions(ctx, logger, input, relationshipDeleteCfg)

		output, err := qb.DeleteRelationship(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	relationshipDeleteCfg, flags = cliutil.AddCommand(relationshipCmd, relationshipDeleteCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeleteRelationshipInput{})
}
