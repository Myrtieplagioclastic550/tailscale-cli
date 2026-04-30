package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the CLI version, set at build time via ldflags.
	Version = "dev"
	// Commit is the git commit hash, set at build time via ldflags.
	Commit = "none"
	// Date is the build date, set at build time via ldflags.
	Date = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Affiche la version de tailscale-cli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("tailscale-cli version %s (commit: %s, built: %s)\n", Version, Commit, Date)
		},
	}
}
