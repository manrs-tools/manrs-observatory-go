package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

func newFavoritesCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "favorites",
		Short: "Interact with favorites",
		Run: func(cmd *cobra.Command, _ []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal().Err(err).Msg("Failed to display help")
			}
		},
	}

	cmd.AddCommand(newFavoritesFindCMD())
	cmd.AddCommand(newFavoritesGetCMD())
	cmd.AddCommand(newFavoritesDeleteCMD())

	return cmd
}

func newFavoritesFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find favorites",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.FavoriteFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			favorites, err := apiClient.FindFavorites(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find favorites")
			}

			if len(favorites) == 0 {
				log.Warn().Msg("No favorites found")
				return
			}

			printToConsole(favorites)
		},
	}

	return cmd
}

func newFavoritesGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get favorite",
		Example: "  client favorites get :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			favorite, err := apiClient.FindFavorite(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get favorite " + args[0])
			}

			printToConsole(favorite)
		},
	}

	return cmd
}

func newFavoritesDeleteCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete favorite",
		Example: "  client favorites delete :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			requireConfirm(cmd, args[0])

			if err := apiClient.DeleteFavorite(cmd.Context(), args[0]); err != nil {
				log.Fatal().Err(err).Msg("Failed to delete favorite")
			}

			log.Info().Msg("Successfully deleted favorite!")
		},
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")

	return cmd
}
