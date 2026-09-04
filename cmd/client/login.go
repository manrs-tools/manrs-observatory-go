package main

import (
	manrs "github.com/manrs-tools/manrs-observatory-go"
	"github.com/spf13/cobra"
)

func newLoginCMD() *cobra.Command {
	var email, password string

	cmd := &cobra.Command{
		Use:     "login",
		Short:   "Login to the API using an existing API key, or a user's email and password.\nYour email and API key are stored locally for convenience; your password is never persisted.",
		Example: "  client login 2fbf890c-aa01-4f5b-907f-a9c2ba634ec2",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				apiKey := args[0]

				apiClient = manrs.New(apiKey, manrs.WithDebug(debug), manrs.WithTimeout(timeout))

				user, err := apiClient.FindSessionUser(cmd.Context())
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to retrieve user information")
				}

				cfg.APIKey = apiKey
				cfg.UserRole = user.Role

				if err = cfg.Save(); err != nil {
					log.Fatal().Err(err).Msg("Could not save config")
				}

				printToConsole("Successfully set api key as " + apiKey)
			} else {
				if email == "" {
					email = cfg.UserEmail
				}

				if email == "" || password == "" {
					log.Fatal().Msg("No key, email or password provided")
				}

				apiClient = manrs.New(cfg.APIKey, manrs.WithDebug(debug), manrs.WithTimeout(timeout))

				apiKey, err := apiClient.Login(cmd.Context(), email, password)
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to login user")
				}

				apiClient.SetAPIKey(apiKey.Key)

				user, err := apiClient.FindSessionUser(cmd.Context())
				if err != nil {
					log.Fatal().Err(err).Msg("Unable to retrieve user information")
				}

				cfg.APIKey = apiKey.Key
				cfg.UserEmail = email
				cfg.UserRole = user.Role

				if err = cfg.Save(); err != nil {
					log.Fatal().Err(err).Msg("Could not save config")
				}

				printToConsole("Successfully logged in!")
				printToConsole("API key set!")
			}
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Login with your email")
	cmd.Flags().StringVar(&password, "password", "", "Login with your password")

	return cmd
}
