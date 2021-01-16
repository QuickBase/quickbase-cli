package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fieldCreateCfg *viper.Viper

var fieldCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a field in a table",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fieldCreateCfg)
			qbcli.SetOptionFromArg(fieldCreateCfg, args, 0, qbclient.OptionTableID)
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateFieldInput{Properties: &qbclient.CreateFieldInputProperties{}}
		qbcli.GetOptions(ctx, logger, input, fieldCreateCfg)

		output, err := qb.CreateField(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fieldCreateCfg, flags = cliutil.AddCommand(fieldCmd, fieldCreateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CreateFieldInput{Properties: &qbclient.CreateFieldInputProperties{}})
}
