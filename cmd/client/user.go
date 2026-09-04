package main

import (
	"github.com/manrs-tools/manrs-observatory-go/model"
	"github.com/spf13/cobra"
)

func newUserCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Show your current user information",
		Run: func(cmd *cobra.Command, _ []string) {
			user, err := apiClient.FindSessionUser(cmd.Context())
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get user info")
			}

			printToConsole(user)
		},
	}

	return cmd
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Interact with users",
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal().Err(err).Msg("Failed to display help")
		}
	},
}

func newUsersCMD() *cobra.Command {
	usersCmd.AddCommand(newUsersFindCMD())
	usersCmd.AddCommand(newUsersGetCMD())
	usersCmd.AddCommand(newUsersUpdateCMD())
	return usersCmd
}

func newUsersFindCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "Find users",
		Run: func(cmd *cobra.Command, _ []string) {
			var filter model.UserFilter

			if err := unmarshalFlags(cmd, &filter); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			users, err := apiClient.FindUsers(cmd.Context(), &filter)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to find users")
			}

			if len(users) == 0 {
				log.Warn().Msg("No users found")
				return
			}

			printToConsole(users)
		},
	}

	cmd.Flags().String("email", "", "Filter by email")
	cmd.Flags().String("firstName", "", "Filter by first name")
	cmd.Flags().String("lastName", "", "Filter by last name")
	cmd.Flags().Bool("mailConformanceEnabled", false, "Only return users with conformance report emails enabled")
	cmd.Flags().Bool("mailListEnabled", false, "Only return users with mailing list emails enabled")
	cmd.Flags().String("participantID", "", "Filter by participant ID")
	cmd.Flags().String("role", "", "Filter by role")
	cmd.Flags().String("status", "", "Filter by status")

	return cmd
}

func newUsersGetCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get user",
		Example: "  client users get :id",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			user, err := apiClient.FindUser(cmd.Context(), args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to get user " + args[0])
			}

			printToConsole(user)
		},
	}

	return cmd
}

func newUsersUpdateCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update user",
		Example: "  client users update :id --email=dev@example.com",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var update model.UserUpdate

			if err := unmarshalFlags(cmd, &update); err != nil {
				log.Fatal().Err(err).Msg("Failed to unmarshal flags")
			}

			user, err := apiClient.UpdateUser(cmd.Context(), args[0], &update)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to update user")
			}

			printToConsole(user)
		},
	}

	cmd.Flags().String("email", "", "Update user's email")
	cmd.Flags().String("firstName", "", "Update user's first name")
	cmd.Flags().String("lastName", "", "Update user's last name")
	cmd.Flags().String("role", "", "Update user's role")
	cmd.Flags().String("status", "", "Update user's status")

	return cmd
}
