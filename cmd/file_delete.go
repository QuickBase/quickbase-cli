package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileDeleteCfg *viper.Viper

var fileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a file from a record",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fileDeleteCfg)
			qbcli.SetOptionFromArg(fileDeleteCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(fileDeleteCfg, args, 1, "field-id")
			qbcli.SetOptionFromArg(fileDeleteCfg, args, 2, "record-id")
			qbcli.SetOptionFromArg(fileDeleteCfg, args, 3, "version")
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.DeleteFileInput{}
		qbcli.GetOptions(ctx, logger, input, fileDeleteCfg)

		output, err := qb.DeleteFile(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fileDeleteCfg, flags = cliutil.AddCommand(fileCmd, fileDeleteCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.DeleteFileInput{})
}
