package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	manrs "github.com/manrs-tools/manrs-observatory-go"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	defaultTimeout = 60 * time.Second
)

var (
	apiClient                     *manrs.Client
	cfg                           *Config
	log                           zerolog.Logger
	debug, prettyLog, writeToFile bool
	timeout                       time.Duration
	format, logLevel              string
	limit                         uint64
	slash                         = string(os.PathSeparator)
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "MANRS Observatory Client",
	Long:  `Interact with the MANRS Observatory API`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		logger, err := newLogger(logLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to setup logger: %s: %v\n", logLevel, err)
			os.Exit(1)
		}

		log = logger

		configDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to retrieve user's home directory")
		}

		cfg, err = newConfig(fmt.Sprintf("%s%s.config%smanrs-observatory-client", strings.TrimSuffix(configDir, slash), slash, slash))
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to initialize config")
		}

		apiClient = manrs.New(cfg.APIKey, manrs.WithDebug(debug), manrs.WithTimeout(timeout))
	},
	Version: manrs.Version,
}

func main() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable console debugging")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "yaml", "Set the output format for CLI commands (json, jsonp, yaml)")
	rootCmd.PersistentFlags().Uint64VarP(&limit, "limit", "l", 0, "Limit the quantity of returned results")
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "info", "Set log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().BoolVar(&prettyLog, "prettyLog", true, "Pretty print logs to console")
	rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", defaultTimeout, "Specify the API HTTP timeout")
	rootCmd.PersistentFlags().BoolVarP(&writeToFile, "writetofile", "w", false, "Write the output to a file")

	rootCmd.AddCommand(newAPIKeysCMD())
	rootCmd.AddCommand(newASCMD())
	rootCmd.AddCommand(newConfigCMD())
	rootCmd.AddCommand(newDocsCMD())
	rootCmd.AddCommand(newEmailsCMD())
	rootCmd.AddCommand(newFavoritesCMD())
	rootCmd.AddCommand(newIncidentFeedbackCMD())
	rootCmd.AddCommand(newLoginCMD())
	rootCmd.AddCommand(newParticipantsCMD())
	rootCmd.AddCommand(newROACMD())
	rootCmd.AddCommand(newStatsCMD())
	rootCmd.AddCommand(newTestimonialsCMD())
	rootCmd.AddCommand(newTicketRepliesCMD())
	rootCmd.AddCommand(newTicketsCMD())
	rootCmd.AddCommand(newUserCMD())
	rootCmd.AddCommand(newUsersCMD())
	rootCmd.AddCommand(newVersionCMD())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func markFlagsRequired(cmd *cobra.Command, flags ...string) error {
	for _, flag := range flags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			return fmt.Errorf("mark required flag %s", flag)
		}
	}

	return nil
}

func marshal(data any) []byte {
	var err error
	var output []byte

	switch strings.ToLower(format) {
	case "json":
		output, err = json.Marshal(data)
	case "jsonp":
		output, err = json.MarshalIndent(data, "", "\t")
	default:
		output, err = yaml.Marshal(data)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to marshal data")
	}

	return output
}

func newLogger(level string) (zerolog.Logger, error) {
	var logger zerolog.Logger

	logLevelParsed, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return logger, fmt.Errorf("invalid log level %s", level)
	}

	var logWriter io.Writer

	if prettyLog {
		logWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		logWriter = os.Stdout
	}

	logger = zerolog.New(logWriter).With().Timestamp().Logger().Level(logLevelParsed)

	return logger, nil
}

func printToConsole(data any) {
	if writeToFile {
		extension := format
		if extension == "jsonp" {
			extension = "json"
		}

		filename := cast.ToString(time.Now().UnixNano()) + "." + extension

		if err := printToFile(data, filename); err != nil {
			log.Fatal().Err(err).Msg("Failed to write output to file")
		}

		log.Info().Msg("Output written to " + filename)

		return
	}

	marshalledData := marshal(data)

	print(string(marshalledData)) //nolint:forbidigo // Print is intentional for raw output without a trailing newline.
}

func printToFile(data any, file string) error {
	marshalledData := marshal(data)

	outputFile, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer outputFile.Close()

	if _, err = outputFile.Write(marshalledData); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func requireConfirm(cmd *cobra.Command, id string) {
	force, _ := cmd.Flags().GetBool("force")
	if force {
		return
	}

	fmt.Fprintf(os.Stderr, "Delete %s? This cannot be undone. Type \"yes\" to confirm: ", id)

	var response string
	_, _ = fmt.Fscanln(os.Stdin, &response)

	if strings.ToLower(response) != "yes" {
		log.Fatal().Msg("Aborted")
	}
}
