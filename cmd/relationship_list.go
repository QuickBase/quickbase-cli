package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relationshipListCfg *viper.Viper

var relationshipListCmd = &cobra.Command{
	Use:   "list",
	Short: "List a table's relationships",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableIDAs(relationshipListCfg, "child-table-id")
			qbcli.SetOptionFromArg(relationshipListCfg, args, 0, "child-table-id")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.ListRelationshipsInput{}
		qbcli.GetOptions(ctx, logger, input, relationshipListCfg)

		output, err := qb.ListRelationships(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	relationshipListCfg, flags = cliutil.AddCommand(relationshipCmd, relationshipListCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.ListRelationshipsInput{})
}
