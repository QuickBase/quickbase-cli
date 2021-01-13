package qbclient

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// EnvPrefix is the prefix for environment variables containing configuration.
const EnvPrefix = "QUICKBASE"

// ConfigFilename is the name of the configuration file.
const ConfigFilename = "config.yml"

// Option* constants contain CLI options.
const (
	OptionAppID          = "app-id"
	OptionConfigDir      = "config-dir"
	OptionFieldID        = "field-id"
	OptionProfile        = "profile"
	OptionRealmHostname  = "realm-hostname"
	OptionRelationshipID = "relationship-id"
	OptionTableID        = "table-id"
	OptionTemporaryToken = "temp-token"
	OptionUserToken      = "user-token"
)

// ConfigIface is implemented by structs used to configure the cleint.
type ConfigIface interface {

	// ConfigDir returns the configuration directory.
	ConfigDir() string

	// DefaultAppID returns the default app ID.
	DefaultAppID() string

	// Default FieldID returns the default field ID
	DefaultFieldID() int

	// DefaultTableID returns the default table ID.
	DefaultTableID() string

	// Profile returns the configured profile.
	Profile() string

	// RealmHostname returns the configured realm hostname.
	RealmHostname() string

	// TemporaryToken returns the configured log level.
	TemporaryToken() string

	// UserToken returns the configured log level.
	UserToken() string
}

// Config contains configuration for the client.
type Config struct {
	cfg *viper.Viper
}

// NewConfig returns a new config
func NewConfig(cfg *viper.Viper) Config {
	return Config{cfg: cfg}
}

// ConfigDir returns the configuration directory.
func (c Config) ConfigDir() string { return c.cfg.GetString(OptionConfigDir) }

// DefaultAppID returns the default app ID.
func (c Config) DefaultAppID() string { return c.cfg.GetString(OptionAppID) }

// DefaultFieldID returns the default field ID.
func (c Config) DefaultFieldID() int { return c.cfg.GetInt(OptionFieldID) }

// DefaultTableID returns the default table ID.
func (c Config) DefaultTableID() string { return c.cfg.GetString(OptionTableID) }

// Profile returns the configured profile.
func (c Config) Profile() string { return c.cfg.GetString(OptionProfile) }

// RealmHostname returns the configured realm hostname.
func (c Config) RealmHostname() string { return c.cfg.GetString(OptionRealmHostname) }

// TemporaryToken returns the configured log level.
func (c Config) TemporaryToken() string { return c.cfg.GetString(OptionTemporaryToken) }

// UserToken returns the configured log level.
func (c Config) UserToken() string { return c.cfg.GetString(OptionUserToken) }

// ReadInConfig reads in configuration from the config file.
func ReadInConfig(cfg *viper.Viper) error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	// Set the default profile and configuration file directory.
	cfg.SetDefault(OptionProfile, "default")
	cfg.SetDefault(OptionConfigDir, Filepath(homeDir, ".config", "quickbase"))

	// Read in configuration from environment variables.
	cfg.SetEnvPrefix(EnvPrefix)
	cfg.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	cfg.AutomaticEnv()

	// Read the configuration file in the configuration directory if it exists.
	configFile, err := ReadConfigFile(cfg.GetString(OptionConfigDir))
	if err != nil {
		return err
	}

	// Get the profile's configuration if set.
	p := cfg.GetString(OptionProfile)
	if config, ok := configFile[p]; ok {
		cfg.SetDefault(OptionRealmHostname, config.RealmHostname)
		cfg.SetDefault(OptionUserToken, config.UserToken)
		cfg.SetDefault(OptionTemporaryToken, config.TemporaryToken)
		cfg.SetDefault(OptionAppID, config.AppID)
		cfg.SetDefault(OptionTableID, config.TableID)
		cfg.SetDefault(OptionFieldID, config.FieldID)
	}

	return nil
}

// ReadConfigFile reads and parses the configuration file.
func ReadConfigFile(dir string) (cf ConfigFile, err error) {
	cf = make(map[string]*ConfigFileProfile, 0)

	filepath := Filepath(dir, ConfigFilename)
	if !FileExists(filepath) {
		return
	}

	var b []byte
	if b, err = ioutil.ReadFile(filepath); err != nil {
		return
	}

	err = yaml.Unmarshal(b, &cf)
	return
}

// WriteConfigFile writes a configuration file.
func WriteConfigFile(dir string, cf ConfigFile) (err error) {

	if !DirExists(dir) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}
	}

	var b []byte
	if b, err = yaml.Marshal(cf); err != nil {
		return
	}

	filepath := Filepath(dir, ConfigFilename)
	err = ioutil.WriteFile(filepath, b, 0600)
	return
}

// ConfigFile models the configuration file.
type ConfigFile map[string]*ConfigFileProfile

// ConfigFileProfile models the configuration for a profile.
type ConfigFileProfile struct {
	RealmHostname  string `yaml:"realm_hostname,omitempty" json:"realm_hostname,omitempty"`
	UserToken      string `yaml:"user_token,omitempty" json:"user_token,omitempty"`
	TemporaryToken string `yaml:"temp_token,omitempty" json:"temp_token,omitempty"`
	AppID          string `yaml:"app_id,omitempty" json:"app_id,omitempty"`
	TableID        string `yaml:"table_id,omitempty" json:"table_id,omitempty"`
	FieldID        int    `yaml:"field_id,omitempty" json:"field_id,omitempty"`
}
