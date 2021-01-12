package qbcli

import (
	"context"
	"errors"
	"fmt"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/cpliakas/cliutil"
	"github.com/go-playground/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Option* constants contain CLI options.
const (
	OptionDumpDirectory  = "dump-dir"
	OptionJMESPathFilter = "filter"
	OptionLogFile        = "log-file"
	OptionLogLevel       = "log-level"
	OptionQuiet          = "quiet"
)

// Option*Description constants contain common option descriptions.
const (
	OptionAppIDDescription         = "unique identifier of an app (required)"
	OptionFieldIDDescription       = "unique identifier (fid) of the field (required)"
	OptionParentTableIDDescription = "unique identifier (dbid) of the parent table (required)"
	OptionTableIDDescription       = "unique identifier (dbid) of the table (required)"
	OptionQuietDescription         = "suppress output written to stdout"
)

// NewGlobalConfig returns a GlobalConfig.
func NewGlobalConfig(cmd *cobra.Command, cfg *viper.Viper) GlobalConfig {
	flags := cliutil.NewFlagger(cmd, cfg)

	flags.PersistentString(OptionDumpDirectory, "d", "", "directory for files that request/response are dumped to for debugging")
	flags.PersistentString(OptionJMESPathFilter, "F", "", "JMESPath filter applied to output")
	flags.PersistentString(OptionLogFile, "f", "", "file log messages are written to")
	flags.PersistentString(OptionLogLevel, "l", cliutil.LogNotice, "minimum log level")
	flags.PersistentString(qbclient.OptionProfile, "p", "default", "configuration profile")
	flags.PersistentBool(OptionQuiet, "q", false, OptionQuietDescription)
	flags.PersistentString(qbclient.OptionRealmHostname, "r", "", "realm hostname, e.g., example.quickbase.com")
	flags.PersistentString(qbclient.OptionTemporaryToken, "t", "", "temporary token used to authenticate API requests")
	flags.PersistentString(qbclient.OptionUserToken, "u", "", "user token used to authenticate API requests")

	return GlobalConfig{cfg: cfg}
}

// GlobalConfig contains configuration common to all commands.
type GlobalConfig struct {
	cfg *viper.Viper
}

// ConfigDir returns the configuration directory.
func (c GlobalConfig) ConfigDir() string { return c.cfg.GetString(qbclient.OptionConfigDir) }

// DefaultAppID returns the default app ID.
func (c GlobalConfig) DefaultAppID() string { return c.cfg.GetString(qbclient.OptionAppID) }

// DefaultFieldID returns the default field ID.
func (c GlobalConfig) DefaultFieldID() int { return c.cfg.GetInt(qbclient.OptionFieldID) }

// DefaultTableID returns the default table ID.
func (c GlobalConfig) DefaultTableID() string { return c.cfg.GetString(qbclient.OptionTableID) }

// DumpDirectory returns the configured dump file directory.
func (c GlobalConfig) DumpDirectory() string { return c.cfg.GetString(OptionDumpDirectory) }

// JMESPathFilter returns the JMESPath filter.
func (c GlobalConfig) JMESPathFilter() string { return c.cfg.GetString(OptionJMESPathFilter) }

// LogFile returns the configured log file.
func (c GlobalConfig) LogFile() string { return c.cfg.GetString(OptionLogFile) }

// LogLevel returns the configured log level.
func (c GlobalConfig) LogLevel() string { return c.cfg.GetString(OptionLogLevel) }

// Profile returns the configured profile.
func (c GlobalConfig) Profile() string { return c.cfg.GetString(qbclient.OptionProfile) }

// Quiet returns whehter to suppress output written to stdout.
func (c GlobalConfig) Quiet() bool { return c.cfg.GetBool(OptionQuiet) }

// RealmHostname returns the configured realm hostname.
func (c GlobalConfig) RealmHostname() string { return c.cfg.GetString(qbclient.OptionRealmHostname) }

// TemporaryToken returns the configured log level.
func (c GlobalConfig) TemporaryToken() string { return c.cfg.GetString(qbclient.OptionTemporaryToken) }

// UserToken returns the configured log level.
func (c GlobalConfig) UserToken() string { return c.cfg.GetString(qbclient.OptionUserToken) }

// ReadInConfig reads in the config file.
func (c *GlobalConfig) ReadInConfig() error { return qbclient.ReadInConfig(c.cfg) }

// Validate reads the configuration file and validates the global configuration
// options.
func (c *GlobalConfig) Validate() error {
	if !cliutil.LogLevelValid(c.LogLevel()) {
		return fmt.Errorf("value %q for option %q: %w", c.LogLevel(), OptionLogLevel, errors.New("invalid value"))
	}

	if err := c.ReadInConfig(); err != nil {
		return err
	}

	if c.RealmHostname() == "" {
		return fmt.Errorf("option %q: %w", qbclient.OptionRealmHostname, errors.New("value required"))
	}

	return nil
}

// SetDefaultAppID sets the default app in the command's configuration.
func (c GlobalConfig) SetDefaultAppID(cfg *viper.Viper) {
	if appID := c.DefaultAppID(); appID != "" {
		cfg.SetDefault(qbclient.OptionAppID, appID)
	}
}

// SetDefaultTableID sets the default table in the command's configuration.
func (c GlobalConfig) SetDefaultTableID(cfg *viper.Viper) {
	if tableID := c.DefaultTableID(); tableID != "" {
		cfg.SetDefault(qbclient.OptionTableID, tableID)
	}
}

// SetOptionFromArg sets an option from an argument.
func SetOptionFromArg(cfg *viper.Viper, args []string, idx int, option string) {
	if len(args) > idx {
		cfg.SetDefault(option, args[idx])
	}
}

// GetOptions gets options based on the input and validates them.
func GetOptions(ctx context.Context, logger *cliutil.LeveledLogger, input interface{}, cfg *viper.Viper) {
	err := cliutil.GetOptions(input, cfg)
	logger.FatalIfError(ctx, "error getting options", err)

	err = validator.New().Struct(input)
	HandleError(ctx, logger, "input not valid", err)
}
