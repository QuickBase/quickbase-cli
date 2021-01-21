package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldCreateSummaryCfg *viper.Viper

var fieldCreateSummaryCmd = &cobra.Command{
	Use:   "create-summary",
	Short: "Create a summary field in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableIDAs(fieldCreateSummaryCfg, "child-table-id")
			qbcli.SetOptionFromArg(fieldCreateSummaryCfg, args, 0, "child-table-id")
			qbcli.SetOptionFromArg(fieldCreateSummaryCfg, args, 1, qbclient.OptionRelationshipID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.UpdateRelationshipInput{}
		qbcli.GetOptions(ctx, logger, input, fieldCreateSummaryCfg)

		sf := &qbclient.RelationshipSummaryField{}
		qbcli.GetOptions(ctx, logger, sf, fieldCreateSummaryCfg)
		input.SummaryFields = []*qbclient.RelationshipSummaryField{sf}

		output, err := qb.UpdateRelationship(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldCreateSummaryCfg, flags = cliutil.AddCommand(fieldCmd, fieldCreateSummaryCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.RelationshipSummaryField{})
	flags.SetOptions(&qbclient.UpdateRelationshipInput{})
}
