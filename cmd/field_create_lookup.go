package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldCreateLookupCfg *viper.Viper

var fieldCreateLookupCmd = &cobra.Command{
	Use:   "create-lookup",
	Short: "Create a lookup field in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableIDs(fieldCreateLookupCfg, "child-table-id")
			qbcli.SetOptionFromArg(fieldCreateLookupCfg, args, 0, "child-table-id")
			qbcli.SetOptionFromArg(fieldCreateLookupCfg, args, 1, qbclient.OptionRelationshipID)
			qbcli.SetOptionFromArg(fieldCreateLookupCfg, args, 2, "lookup-field-ids")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateRelationshipInput{}
		qbcli.GetOptions(ctx, logger, input, fieldCreateLookupCfg)

		output, err := qb.UpdateRelationship(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldCreateLookupCfg, flags = cliutil.AddCommand(fieldCmd, fieldCreateLookupCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.UpdateRelationshipInput{})
}
