package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Interact with statistics",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newStatsCMD() *cobra.Command {
	statsCmd.AddCommand(newStatsConformanceCMD())
	statsCmd.AddCommand(newStatsSummaryCMD())
	statsCmd.AddCommand(newStatsHistoryCMD())
	statsCmd.AddCommand(newStatsOverviewCMD())
	return statsCmd
}

func newStatsConformanceCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conformance",
		Short: "Get participant conformance scores",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.ConformanceFilter
			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			scores, err := apiClient.GetConformanceScores(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get conformance scores")
			}

			printToConsole(scores)
		},
	}

	cmd.Flags().String("type", "", "Filter by programme type (netops, ixp, cdn, vendor)")

	return cmd
}

func newStatsSummaryCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Get the aggregate score summary",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.ScoresSummaryFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			summary, err := apiClient.GetScoresSummary(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get score summary")
			}

			printToConsole(summary)
		},
	}

	cmd.Flags().String("asns", "", "Filter by ASNs (comma-separated)")
	cmd.Flags().String("economies", "", "Filter by economies (comma-separated ISO-2 codes)")
	cmd.Flags().String("month", "", "Report month (RFC3339, e.g. 2026-05-01T00:00:00Z)")
	cmd.Flags().String("region", "", "Filter by region")

	return cmd
}

func newStatsHistoryCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Get statistics history",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.StatsHistoryFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			history, err := apiClient.GetStatsHistory(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get stats history")
			}

			printToConsole(history)
		},
	}

	cmd.Flags().String("asn", "", "Filter by ASN")
	cmd.Flags().String("region", "", "Filter by region")

	return cmd
}

func newStatsOverviewCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "Get statistics overview",
		Run: func(cmd *cobra.Command, _ []string) {
			overview, err := apiClient.GetStatsOverview(cmd.Context())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get stats overview")
			}

			printToConsole(overview)
		},
	}

	return cmd
}
