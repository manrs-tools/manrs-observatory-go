package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

var participantsCmd = &cobra.Command{
	Use:   "participants",
	Short: "Interact with participants",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newParticipantsCMD() *cobra.Command {
	participantsCmd.AddCommand(newParticipantsFindCMD())
	participantsCmd.AddCommand(newParticipantsGetCMD())
	return participantsCmd
}

func newParticipantsFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find participants",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.ParticipantFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			participants, err := apiClient.FindParticipants(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find participants")
			}

			if len(participants) == 0 {
				log.Warn().Msg("No participants found")
				return
			}

			printToConsole(participants)
		},
	}

	cmd.Flags().String("abbreviation", "", "Filter by abbreviation")
	cmd.Flags().String("asn", "", "Filter by ASN")
	cmd.Flags().Bool("deleted", false, "Include deleted participants")
	cmd.Flags().String("ids", "", "Filter by IDs (comma-separated)")
	cmd.Flags().String("legacyID", "", "Filter by legacy ID")
	cmd.Flags().String("name", "", "Filter by name")
	cmd.Flags().String("primaryUserID", "", "Filter by primary user ID")
	cmd.Flags().String("region", "", "Filter by region")
	cmd.Flags().String("status", "", "Filter by status")
	cmd.Flags().String("type", "", "Filter by type")

	return cmd
}

func newParticipantsGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get participant",
		Example: "  client participants get :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			participant, err := apiClient.FindParticipant(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get participant " + args[0])
			}

			printToConsole(participant)
		},
	}

	return cmd
}
