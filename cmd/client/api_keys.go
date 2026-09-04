package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newAPIKeysCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apiKeys",
		Aliases: []string{"apikeys"},
		Short:   "Interact with API keys",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal().Err(err).Msg("Failed to display help")
			}
		},
	}

	cmd.AddCommand(newAPIKeysCreateCMD())
	cmd.AddCommand(newAPIKeysDeleteCMD())
	cmd.AddCommand(newAPIKeysFindCMD())
	cmd.AddCommand(newAPIKeysGetCMD())

	return cmd
}

func newAPIKeysCreateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create API key",
		Example: "  client apikeys create --description=mykey",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, _ []string) {
			var apiKey model.APIKey

			cmd.Flags().Visit(func(flag *pflag.Flag) {
				if flag.Name == "description" {
					apiKey.Description = flag.Value.String()
				}
			})

			created, err := apiClient.CreateAPIKey(cmd.Context(), &apiKey)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to create API key")
			}

			printToConsole(created)
		},
	}

	cmd.Flags().String("description", "", "Set API key description")

	return cmd
}

func newAPIKeysDeleteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete API key",
		Example: "  client apikeys delete c318fc9b-d343-4a8f-b3cb-1a18773ddf05",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			requireConfirm(cmd, args[0])

			if err := apiClient.DeleteAPIKey(cmd.Context(), args[0]); err != nil {
				log.Fatal().Err(err).Msg("Failed to delete API key")
			}

			log.Info().Msg("Successfully deleted API key!")
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")

	return cmd
}

func newAPIKeysFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find API keys",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.APIKeyFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			apiKeys, err := apiClient.FindAPIKeys(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find API keys")
			}

			if len(apiKeys) == 0 {
				log.Warn().Msg("No API keys found")
				return
			}

			printToConsole(apiKeys)
		},
	}

	cmd.Flags().String("userID", "", "Filter by user ID")

	return cmd
}

func newAPIKeysGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get API key",
		Example: "  client apikeys get 19ad4d0e-569d-4ee7-9ee5-c24594a1acd9",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			apiKey, err := apiClient.FindAPIKey(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get API key " + args[0])
			}

			printToConsole(apiKey)
		},
	}

	return cmd
}
