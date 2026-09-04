package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var ticketRepliesCmd = &cobra.Command{
	Use:   "ticket-replies",
	Short: "Interact with ticket replies",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newTicketRepliesCMD() *cobra.Command {
	ticketRepliesCmd.AddCommand(newTicketRepliesCreateCMD())
	return ticketRepliesCmd
}

func newTicketRepliesCreateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create ticket reply",
		Example: "  client ticket-replies create --ticketID=:id --message='...'",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			var reply model.TicketReply

			cmd.Flags().Visit(func(flag *pflag.Flag) {
				switch flag.Name {
				case "ticketID":
					reply.TicketID = flag.Value.String()
				case "message":
					reply.Message = flag.Value.String()
				}
			})

			if err := apiClient.CreateTicketReply(cmd.Context(), &reply); err != nil {
				log.Fatal().Err(err).Msg("Failed to create ticket reply")
			}

			log.Info().Msg("Successfully created ticket reply!")
		},
	}

	cmd.Flags().String("ticketID", "", "Ticket ID")
	cmd.Flags().String("message", "", "Reply message")
	_ = markFlagsRequired(cmd, "ticketID", "message")

	return cmd
}
