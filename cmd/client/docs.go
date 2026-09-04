package main

import (
	"os/exec"
	"runtime"

	manrs "github.com/manrs-tools/manrs-observatory-go"
	"github.com/spf13/cobra"
)

func newDocsCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Access the OpenAPI docs",
		Run: func(cmd *cobra.Command, _ []string) {
			var cmdArgs []string
			var cmdName string

			switch os := runtime.GOOS; os {
			case "darwin":
				cmdName = "open"
			case "windows":
				cmdArgs = []string{"/c", "start"}
				cmdName = "cmd"
			default:
				cmdName = "xdg-open"
			}

			cmdArgs = append(cmdArgs, manrs.DocsURL)

			execCMD := exec.CommandContext(cmd.Context(), cmdName, cmdArgs...)
			if err := execCMD.Start(); err != nil {
				log.Fatal().Err(err).Msg("Failed to open browser to " + manrs.DocsURL)
			}

			log.Info().Msg("Browser opened to " + manrs.DocsURL)
		},
	}

	return cmd
}
