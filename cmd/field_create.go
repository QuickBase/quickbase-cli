package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldCreateCfg *viper.Viper

// fieldCreateCmd represents the app get command
var fieldCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a field in a table",
	Args:  fieldCreateCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateFieldInput{
			TableID: fieldCreateCfg.GetString(qbclient.OptionTableID),
			Field:   qbcli.NewFieldFromOptions(fieldCreateCfg),
			Properties: &qbclient.CreateFieldInputProperties{
				FieldProperties: qbcli.NewPropertiesFromOptions(fieldCreateCfg),
			},
		}

		if !fieldCreateCfg.GetBool("dry-run") {
			output, err := qb.CreateField(input)
			qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
		} else {
			cliutil.PrintJSON(input)
		}
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldCreateCfg, flags = cliutil.AddCommand(fieldCmd, fieldCreateCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionTableID, "", "", qbcli.OptionTableIDDescription)
	flags.Bool("dry-run", "", false, "print input but don't send it")

	qbcli.SetFieldOptions(flags)
}

func fieldCreateCmdValidate(cmd *cobra.Command, args []string) (err error) {
	if err = globalCfg.Validate(); err == nil {
		globalCfg.SetDefaultTableID(fieldCreateCfg)
		qbcli.SetOptionFromArg(fieldCreateCfg, args, 0, qbclient.OptionTableID)
	}
	return
}
