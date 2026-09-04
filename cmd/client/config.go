package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConfigCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure your MANRS Observatory client",
	}

	cmd.AddCommand(newConfigGetCMD())
	cmd.AddCommand(newConfigSetCMD())
	cmd.AddCommand(newConfigShowCMD())

	return cmd
}

func newConfigGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a config value",
		Example: "  client config get apikey",
		Args:    cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			switch strings.ToLower(args[0]) {
			case "apikey":
				printToConsole("api key: " + cfg.APIKey)
			case "useremail":
				printToConsole("user email: " + cfg.UserEmail)
			case "userrole":
				printToConsole("user role: " + cfg.UserRole)
			default:
				log.Fatal().Msg("Unknown config key")
			}
		},
	}

	return cmd
}

func newConfigSetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Short:   "Set a config value",
		Example: "  client config set apikey 019a0dd4-11a5-7477-91a8-538b1bc334e4",
		Args:    cobra.ExactArgs(2), //nolint:mnd // Unnecessary.
		Run: func(_ *cobra.Command, args []string) {
			switch strings.ToLower(args[0]) {
			case "apikey":
				cfg.APIKey = args[1]
			case "useremail":
				cfg.UserEmail = args[1]
			case "userrole":
				cfg.UserRole = args[1]
			default:
				log.Fatal().Msg("Unknown config key")
			}

			if err := cfg.Save(); err != nil {
				log.Fatal().Err(err).Msg("Unable to save config")
			}

			log.Info().Msg("Config updated")
		},
	}

	return cmd
}

func newConfigShowCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Print full config",
		Example: "  client config show",
		Args:    cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			printToConsole(struct {
				APIKey    string `json:"apiKey"    yaml:"apiKey"`
				UserEmail string `json:"userEmail" yaml:"userEmail"`
				UserRole  string `json:"userRole"  yaml:"userRole"`
			}{
				APIKey:    "[redacted]",
				UserEmail: cfg.UserEmail,
				UserRole:  cfg.UserRole,
			})
		},
	}

	return cmd
}

// Config holds client configuration persisted to disk.
type Config struct {
	dir       string
	path      string
	APIKey    string `json:"apiKey"    yaml:"apiKey"`
	UserEmail string `json:"userEmail" yaml:"userEmail"`
	UserRole  string `json:"userRole"  yaml:"userRole"`
}

func newConfig(directory string) (*Config, error) {
	config := Config{
		dir:       directory,
		path:      directory + slash + "config.yml",
		APIKey:    "",
		UserEmail: "",
		UserRole:  "",
	}

	if err := config.Load(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &config, nil
}

func (c *Config) Load() error {
	if _, err := os.Stat(c.path); errors.Is(err, fs.ErrNotExist) {
		if err = os.MkdirAll(c.dir, 0o750); err != nil {
			return fmt.Errorf("create config directory: %w", err)
		}

		if err = c.Save(); err != nil {
			return err
		}
	}

	configData, err := os.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	if err = yaml.Unmarshal(configData, &c); err != nil {
		return fmt.Errorf("unmarshal config values: %w", err)
	}

	return nil
}

func (c *Config) Save() error {
	configData, err := yaml.Marshal(c) //nolint:gosec // Persisting the API key is this file's purpose; written 0600 below.
	if err != nil {
		return fmt.Errorf("marshal config values: %w", err)
	}

	if err = os.WriteFile(c.path, configData, 0o600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}
