package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Génère le script d'autocomplétion pour le shell spécifié",
		Long: `Génère le script d'autocomplétion pour tailscale-cli pour le shell spécifié.

Exemples :

  # Bash
  source <(tailscale-cli completion bash)

  # Zsh
  tailscale-cli completion zsh > "${fpath[1]}/_tailscale-cli"

  # Fish
  tailscale-cli completion fish | source

  # PowerShell
  tailscale-cli completion powershell | Out-String | Invoke-Expression
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}

	return completionCmd
}
