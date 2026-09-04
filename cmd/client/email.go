package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

var emailCmd = &cobra.Command{
	Use:   "emails",
	Short: "Interact with emails",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newEmailsCMD() *cobra.Command {
	emailCmd.AddCommand(newEmailsFindSentCMD())
	emailCmd.AddCommand(newEmailsGetSentCMD())
	return emailCmd
}

func newEmailsFindSentCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find sent emails",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.EmailFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			emails, err := apiClient.FindSentEmails(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find emails")
			}

			if len(emails) == 0 {
				log.Warn().Msg("No emails found")
				return
			}

			printToConsole(emails)
		},
	}

	cmd.Flags().String("recipient", "", "Filter by recipient email")
	cmd.Flags().String("userID", "", "Filter by user ID")

	return cmd
}

func newEmailsGetSentCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get sent email",
		Example: "  client emails get :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			email, err := apiClient.FindSentEmail(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get email " + args[0])
			}

			printToConsole(email)
		},
	}

	return cmd
}
