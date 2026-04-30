package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// DnsOptions contains closures for lazy resolution of shared options.
type DnsOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// readBody reads JSON content from a file or stdin.
func readBody(file string, stdin bool) ([]byte, error) {
	if stdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
		return data, nil
	}
	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", file, err)
		}
		return data, nil
	}
	return nil, fmt.Errorf("either --file or --stdin must be specified")
}

// NewCmdDns returns the dns parent command with all its subcommands.
func NewCmdDns(opts DnsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Gestion du DNS Tailscale",
		Long:  "Commandes pour configurer les nameservers, les search paths, les préférences DNS, le split DNS et la configuration DNS du tailnet.",
	}

	cmd.AddCommand(newNameserversCmd(opts))
	cmd.AddCommand(newPreferencesCmd(opts))
	cmd.AddCommand(newSearchpathsCmd(opts))
	cmd.AddCommand(newSplitCmd(opts))
	cmd.AddCommand(newConfigCmd(opts))

	return cmd
}

// --- nameservers ---

func newNameserversCmd(opts DnsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nameservers",
		Short: "Gestion des nameservers DNS",
	}

	cmd.AddCommand(newNameserversListCmd(opts))
	cmd.AddCommand(newNameserversSetCmd(opts))

	return cmd
}

func newNameserversListCmd(opts DnsOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Liste les nameservers du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			path := fmt.Sprintf("/tailnet/%s/dns/nameservers", tailnet)

			data, err := client.Get(path)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}
}

func newNameserversSetCmd(opts DnsOptions) *cobra.Command {
	var nameservers []string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Définit les nameservers du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			body := map[string]interface{}{
				"dns": nameservers,
			}

			path := fmt.Sprintf("/tailnet/%s/dns/nameservers", tailnet)

			data, err := client.DoJSON("POST", path, body)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringSliceVar(&nameservers, "nameservers", nil, "Liste des nameservers (ex: 8.8.8.8,1.1.1.1)")
	_ = cmd.MarkFlagRequired("nameservers")

	return cmd
}

// --- preferences ---

func newPreferencesCmd(opts DnsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preferences",
		Short: "Gestion des préférences DNS",
	}

	cmd.AddCommand(newPreferencesGetCmd(opts))
	cmd.AddCommand(newPreferencesSetCmd(opts))

	return cmd
}

func newPreferencesGetCmd(opts DnsOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Récupère les préférences DNS du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			path := fmt.Sprintf("/tailnet/%s/dns/preferences", tailnet)

			data, err := client.Get(path)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}
}

func newPreferencesSetCmd(opts DnsOptions) *cobra.Command {
	var magicDNS bool

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Définit les préférences DNS du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			body := map[string]interface{}{
				"magicDNS": magicDNS,
			}

			path := fmt.Sprintf("/tailnet/%s/dns/preferences", tailnet)

			data, err := client.DoJSON("POST", path, body)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().BoolVar(&magicDNS, "magic-dns", false, "Activer ou désactiver MagicDNS")

	return cmd
}

// --- searchpaths ---

func newSearchpathsCmd(opts DnsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "searchpaths",
		Short: "Gestion des search paths DNS",
	}

	cmd.AddCommand(newSearchpathsListCmd(opts))
	cmd.AddCommand(newSearchpathsSetCmd(opts))

	return cmd
}

func newSearchpathsListCmd(opts DnsOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Liste les search paths DNS du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			path := fmt.Sprintf("/tailnet/%s/dns/searchpaths", tailnet)

			data, err := client.Get(path)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}
}

func newSearchpathsSetCmd(opts DnsOptions) *cobra.Command {
	var searchPaths []string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Définit les search paths DNS du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			body := map[string]interface{}{
				"searchPaths": searchPaths,
			}

			path := fmt.Sprintf("/tailnet/%s/dns/searchpaths", tailnet)

			data, err := client.DoJSON("POST", path, body)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringSliceVar(&searchPaths, "search-paths", nil, "Liste des search paths (ex: example.com,corp.example.com)")
	_ = cmd.MarkFlagRequired("search-paths")

	return cmd
}

// --- split dns ---

func newSplitCmd(opts DnsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "split",
		Short: "Gestion du split DNS",
	}

	cmd.AddCommand(newSplitGetCmd(opts))
	cmd.AddCommand(newSplitUpdateCmd(opts))
	cmd.AddCommand(newSplitSetCmd(opts))

	return cmd
}

func newSplitGetCmd(opts DnsOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Récupère la configuration split DNS du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			path := fmt.Sprintf("/tailnet/%s/dns/split-dns", tailnet)

			data, err := client.Get(path)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}
}

func newSplitUpdateCmd(opts DnsOptions) *cobra.Command {
	var (
		file    string
		stdin   bool
		domain  string
		servers []string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Met à jour partiellement la configuration split DNS (PATCH)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			var bodyBytes []byte

			if domain != "" && len(servers) > 0 {
				// Mode inline : un seul domaine avec ses serveurs
				body := map[string]interface{}{
					domain: servers,
				}
				var err error
				bodyBytes, err = json.Marshal(body)
				if err != nil {
					return fmt.Errorf("marshaling body: %w", err)
				}
			} else {
				// Mode fichier ou stdin
				var err error
				bodyBytes, err = readBody(file, stdin)
				if err != nil {
					return err
				}
			}

			apiPath := fmt.Sprintf("/tailnet/%s/dns/split-dns", tailnet)

			data, err := client.Patch(apiPath, bytes.NewReader(bodyBytes))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Chemin vers le fichier JSON de split DNS")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Lire la configuration depuis stdin")
	cmd.Flags().StringVar(&domain, "domain", "", "Domaine pour le split DNS (mode inline)")
	cmd.Flags().StringSliceVar(&servers, "servers", nil, "Serveurs DNS pour le domaine (mode inline)")

	return cmd
}

func newSplitSetCmd(opts DnsOptions) *cobra.Command {
	var (
		file  string
		stdin bool
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Remplace entièrement la configuration split DNS (PUT)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			bodyBytes, err := readBody(file, stdin)
			if err != nil {
				return err
			}

			apiPath := fmt.Sprintf("/tailnet/%s/dns/split-dns", tailnet)

			data, err := client.Put(apiPath, bytes.NewReader(bodyBytes))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Chemin vers le fichier JSON de split DNS")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Lire la configuration depuis stdin")

	return cmd
}

// --- config ---

func newConfigCmd(opts DnsOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Gestion de la configuration DNS",
	}

	cmd.AddCommand(newConfigGetCmd(opts))
	cmd.AddCommand(newConfigSetCmd(opts))

	return cmd
}

func newConfigGetCmd(opts DnsOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Récupère la configuration DNS du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			path := fmt.Sprintf("/tailnet/%s/dns/configuration", tailnet)

			data, err := client.Get(path)
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}
}

func newConfigSetCmd(opts DnsOptions) *cobra.Command {
	var (
		file  string
		stdin bool
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Définit la configuration DNS du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			bodyBytes, err := readBody(file, stdin)
			if err != nil {
				return err
			}

			apiPath := fmt.Sprintf("/tailnet/%s/dns/configuration", tailnet)

			data, err := client.Post(apiPath, bytes.NewReader(bodyBytes))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Chemin vers le fichier JSON de configuration DNS")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Lire la configuration depuis stdin")

	return cmd
}
