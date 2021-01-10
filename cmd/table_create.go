package cmd

import (
	"github.com/QuickBase/quickbase-cli/qbcli"
	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tableCreateCfg *viper.Viper

// tableCreateCmd represents the app get command
var tableCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a table",
	Args:  tableCreateCmdValidate,

	Run: func(cmd *cobra.Command, args []string) {
		ctx, logger, qb := qbcli.NewClient(cmd, globalCfg)

		input := &qbclient.CreateTableInput{
			AppID:        tableCreateCfg.GetString(qbclient.OptionAppID),
			Name:         tableCreateCfg.GetString("name"),
			Description:  tableCreateCfg.GetString("description"),
			IconName:     tableCreateCfg.GetString("icon-name"),
			SingularNoun: tableCreateCfg.GetString("singular-noun"),
			PluralNoun:   tableCreateCfg.GetString("plural-noun"),
		}

		output, err := qb.CreateTable(input)
		qbcli.Render(ctx, logger, cmd, globalCfg, output, err)
	},
}

func init() {
	var flags *cliutil.Flagger
	tableCreateCfg, flags = cliutil.AddCommand(tableCmd, tableCreateCmd, qbclient.EnvPrefix)

	flags.String(qbclient.OptionAppID, "", "", qbcli.OptionAppIDDescription)
	flags.String("name", "", "", "name for the table")
	flags.String("description", "", "", "description for the table")
	flags.String("icon-name", "", "", "icon for the table")
	flags.String("singular-noun", "", "", "singular noun for records in the table")
	flags.String("plural-noun", "", "", "plural noun for records in the table")
}

func tableCreateCmdValidate(cmd *cobra.Command, args []string) error {
	if err := globalCfg.Validate(); err != nil {
		return err
	}

	// Set the default App ID if configured.
	if appID := globalCfg.DefaultAppID(); appID != "" {
		tableCreateCfg.SetDefault(qbclient.OptionAppID, appID)
	}

	return nil
}
