package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

var ticketsCmd = &cobra.Command{
	Use:   "tickets",
	Short: "Interact with tickets",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newTicketsCMD() *cobra.Command {
	ticketsCmd.AddCommand(newTicketsFindCMD())
	ticketsCmd.AddCommand(newTicketsGetCMD())
	return ticketsCmd
}

func newTicketsFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find tickets",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.TicketFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			tickets, err := apiClient.FindTickets(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find tickets")
			}

			if len(tickets) == 0 {
				log.Warn().Msg("No tickets found")
				return
			}

			printToConsole(tickets)
		},
	}

	cmd.Flags().String("asns", "", "Filter by ASNs (comma-separated)")
	cmd.Flags().String("assignedTo", "", "Filter by assignee")
	cmd.Flags().String("participantID", "", "Filter by participant ID")
	cmd.Flags().String("queue", "", "Filter by queue")
	cmd.Flags().String("status", "", "Filter by status")
	cmd.Flags().String("statusNot", "", "Filter by excluded statuses (comma-separated)")
	cmd.Flags().String("type", "", "Filter by type")
	cmd.Flags().String("userID", "", "Filter by user ID")

	return cmd
}

func newTicketsGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get ticket",
		Example: "  client tickets get :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ticket, err := apiClient.FindTicket(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get ticket " + args[0])
			}

			printToConsole(ticket)
		},
	}

	return cmd
}
