package cmd

import (
	"path"

	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileCreateCfg *viper.Viper

var fileCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a file in a record",

	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = globalCfg.Validate(); err == nil {
			globalCfg.SetDefaultTableID(fileCreateCfg)
			qbcli.SetOptionFromArg(fileCreateCfg, args, 0, qbclient.OptionTableID)
			qbcli.SetOptionFromArg(fileCreateCfg, args, 1, "field-id")
			qbcli.SetOptionFromArg(fileCreateCfg, args, 2, "record-id")
			fileCreateCfg.SetDefault("file-name", path.Base(fileCreateCfg.GetString("file-data")))
		}
		return
	},

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateFileInput{}
		qbcli.GetOptions(ctx, logger, input, fileCreateCfg)

		file := &qbclient.CreateFileInputField{}
		qbcli.GetOptions(ctx, logger, file, fileCreateCfg)
		input.Fields = []*qbclient.CreateFileInputField{file}

		output, err := qb.CreateFile(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	fileCreateCfg, flags = cliutil.AddCommand(fileCmd, fileCreateCmd, qbclient.EnvPrefix)
	flags.SetOptions(&qbclient.CreateFileInput{})
	flags.SetOptions(&qbclient.CreateFileInputField{})
}
