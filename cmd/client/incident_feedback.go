package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

var incidentFeedbackCmd = &cobra.Command{
	Use:   "feedback",
	Short: "Interact with incident feedback",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newIncidentFeedbackCMD() *cobra.Command {
	incidentFeedbackCmd.AddCommand(newIncidentFeedbackDeleteCMD())
	incidentFeedbackCmd.AddCommand(newIncidentFeedbackFindCMD())
	incidentFeedbackCmd.AddCommand(newIncidentFeedbackGetCMD())
	incidentFeedbackCmd.AddCommand(newIncidentSummaryCMD())
	return incidentFeedbackCmd
}

func newIncidentFeedbackFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find incident feedbacks",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.IncidentFeedbackFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			feedbacks, err := apiClient.FindIncidentFeedbacks(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find incident feedbacks")
			}

			if len(feedbacks) == 0 {
				log.Warn().Msg("No incident feedbacks found")
				return
			}

			printToConsole(feedbacks)
		},
	}

	cmd.Flags().String("asn", "", "Filter by ASN")
	cmd.Flags().String("incidentID", "", "Filter by incident ID")
	cmd.Flags().String("reportMonth", "", "Filter by report month")
	cmd.Flags().String("participantID", "", "Filter by participant ID")

	return cmd
}

func newIncidentFeedbackGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get incident feedback",
		Example: "  client feedback get :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			feedback, err := apiClient.FindIncidentFeedback(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get incident feedback " + args[0])
			}

			printToConsole(feedback)
		},
	}

	return cmd
}

func newIncidentFeedbackDeleteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete incident feedback",
		Example: "  client feedback delete :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			requireConfirm(cmd, args[0])

			if err := apiClient.DeleteIncidentFeedback(cmd.Context(), args[0]); err != nil {
				log.Fatal().Err(err).Msg("Failed to delete incident feedback")
			}

			log.Info().Msg("Successfully deleted incident feedback!")
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")

	return cmd
}

func newIncidentSummaryCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Get incident summary",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.IncidentSummaryFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			rows, err := apiClient.GetIncidentSummary(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get incident summary")
			}

			printToConsole(rows)
		},
	}

	cmd.Flags().String("reportMonth", "", "Report month (RFC3339, e.g. 2026-05-01T00:00:00Z)")

	return cmd
}
