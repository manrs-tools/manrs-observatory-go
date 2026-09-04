package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

func newROACMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roa",
		Short: "Interact with ROA data",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal().Err(err).Msg("Failed to display help")
			}
		},
	}

	cmd.AddCommand(newROAFindCMD())
	cmd.AddCommand(newROAStatsFindCMD())

	return cmd
}

func newROAFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find ROAs",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.ROAFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			roas, err := apiClient.FindROAs(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find ROAs")
			}

			if len(roas) == 0 {
				log.Warn().Msg("No ROAs found")
				return
			}

			printToConsole(roas)
		},
	}

	cmd.Flags().String("asns", "", "Filter by ASNs (comma-separated)")
	cmd.Flags().String("ipPrefix", "", "Filter by IP prefix")
	cmd.Flags().String("ipVersion", "", "Filter by IP version")
	cmd.Flags().String("maxLength", "", "Filter by max length")
	cmd.Flags().String("notAfter", "", "Filter by not-after date (RFC3339)")
	cmd.Flags().String("notBefore", "", "Filter by not-before date (RFC3339)")
	cmd.Flags().String("source", "", "Filter by source")
	cmd.Flags().String("sourceDate", "", "Filter by source date (RFC3339)")
	cmd.Flags().Bool("valid", false, "Filter by valid ROAs only")

	return cmd
}

func newROAStatsFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Find ROA stats",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.ROAStatFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			stats, err := apiClient.FindROAAnnouncementRecords(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find ROA stats")
			}

			if len(stats) == 0 {
				log.Warn().Msg("No ROA stats found")
				return
			}

			printToConsole(stats)
		},
	}

	cmd.Flags().String("asns", "", "Filter by ASNs (comma-separated)")
	cmd.Flags().String("ipPrefix", "", "Filter by IP prefix")
	cmd.Flags().String("ipVersion", "", "Filter by IP version")
	cmd.Flags().String("sourceDate", "", "Filter by source date (RFC3339)")
	cmd.Flags().String("sourceDateFrom", "", "Filter by source date from (RFC3339)")
	cmd.Flags().String("sourceDateTo", "", "Filter by source date to (RFC3339)")
	cmd.Flags().String("state", "", "Filter by state")

	return cmd
}
