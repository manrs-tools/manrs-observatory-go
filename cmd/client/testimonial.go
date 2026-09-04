package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

var testimonialsCmd = &cobra.Command{
	Use:   "testimonials",
	Short: "Interact with testimonials",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newTestimonialsCMD() *cobra.Command {
	testimonialsCmd.AddCommand(newTestimonialsFindCMD())
	return testimonialsCmd
}

func newTestimonialsFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find testimonials",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.TestimonialFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			testimonials, err := apiClient.FindTestimonials(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find testimonials")
			}

			if len(testimonials) == 0 {
				log.Warn().Msg("No testimonials found")
				return
			}

			printToConsole(testimonials)
		},
	}

	cmd.Flags().String("participantID", "", "Filter by participant ID")
	cmd.Flags().String("name", "", "Filter by name")

	return cmd
}
