package main

import (
	"strconv"

	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

var asCmd = &cobra.Command{
	Use:   "as",
	Short: "Interact with autonomous systems",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newASCMD() *cobra.Command {
	asCmd.AddCommand(newASFindCMD())
	asCmd.AddCommand(newASGetCMD())
	asCmd.AddCommand(newASIncidentsCMD())
	asCmd.AddCommand(newASScoresCMD())
	return asCmd
}

func newASFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find autonomous systems",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.ASFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			ases, err := apiClient.FindASes(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find ASes")
			}

			if len(ases) == 0 {
				log.Warn().Msg("No ASes found")
				return
			}

			printToConsole(ases)
		},
	}

	cmd.Flags().String("economies", "", "Filter by economies (comma-separated)")

	return cmd
}

func newASGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get autonomous system",
		Example: "  client as get 1234",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			asn, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				log.Fatal().Err(err).Msg("Invalid ASN")
			}

			as, err := apiClient.FindAS(cmd.Context(), uint32(asn))
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get AS " + args[0])
			}

			printToConsole(as)
		},
	}

	return cmd
}

func newASScoresCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scores",
		Short:   "Get ASN scores",
		Example: "  client as scores 1234",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			asn, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				log.Fatal().Err(err).Msg("Invalid ASN")
			}

			var filter model.ScoresDetailsFilter
			if err = unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}
			filter.ASN = uint32(asn)

			scores, history, err := apiClient.GetASNScores(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get ASN scores")
			}

			printToConsole(struct {
				Score   *model.ASNScoreDetail         `json:"score"`
				History []*model.ASNScoreHistoryPoint `json:"history"`
			}{
				Score:   scores,
				History: history,
			})
		},
	}

	cmd.Flags().String("month", "", "Report month (RFC3339, e.g. 2026-05-01T00:00:00Z)")

	return cmd
}

func newASIncidentsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "incidents",
		Short:   "Get ASN incidents",
		Example: "  client as incidents 1234",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			asn, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				log.Fatal().Err(err).Msg("Invalid ASN")
			}

			var filter model.IncidentsFilter
			if err = unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}
			filter.ASNs = []uint32{uint32(asn)}

			incidents, err := apiClient.GetASNIncidents(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get ASN incidents")
			}

			printToConsole(incidents)
		},
	}

	cmd.Flags().String("month", "", "Report month (RFC3339, e.g. 2026-05-01T00:00:00Z)")

	return cmd
}
