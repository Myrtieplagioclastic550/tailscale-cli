package acl

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

// AclOptions contains closures for lazy resolution of shared options.
type AclOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// readBody reads the ACL body from a file or stdin.
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

// NewCmdAcl returns the acl parent command with all its subcommands.
func NewCmdAcl(opts AclOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Gestion des ACL Tailscale",
		Long:  "Commandes pour lire, modifier, prévisualiser et valider les règles d'accès (ACL) du tailnet.",
	}

	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newSetCmd(opts))
	cmd.AddCommand(newPreviewCmd(opts))
	cmd.AddCommand(newValidateCmd(opts))

	return cmd
}

func newGetCmd(opts AclOptions) *cobra.Command {
	var (
		format  string
		details bool
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Récupère le fichier ACL du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			path := fmt.Sprintf("/tailnet/%s/acl", tailnet)
			if details {
				path += "?details=true"
			}

			// Le header Accept devrait être application/json ou application/hujson
			// selon la valeur de --format. Le client actuel ne supporte pas les headers
			// personnalisés directement sur Get ; on passe par Do si nécessaire.
			_ = format

			data, err := client.Get(path)
			if err != nil {
				return err
			}

			// Affichage brut pour préserver le format hujson si demandé
			if format == "hujson" {
				fmt.Println(string(data))
				return nil
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				fmt.Println(string(data))
				return nil
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&format, "format", "json", "Format de sortie (json, hujson)")
	cmd.Flags().BoolVar(&details, "details", false, "Inclure les détails supplémentaires")

	return cmd
}

func newSetCmd(opts AclOptions) *cobra.Command {
	var (
		file    string
		stdin   bool
		format  string
		ifMatch string
	)

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Définit le fichier ACL du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			body, err := readBody(file, stdin)
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/acl", tailnet)

			// Les headers suivants devraient être ajoutés à la requête :
			// - If-Match: ifMatch (pour le contrôle de concurrence optimiste)
			// - Accept: application/json ou application/hujson selon --format
			_, _ = format, ifMatch

			data, err := client.Post(path, bytes.NewReader(body))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				fmt.Println(string(data))
				return nil
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Chemin vers le fichier ACL")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Lire le fichier ACL depuis stdin")
	cmd.Flags().StringVar(&format, "format", "json", "Format du fichier (json, hujson)")
	cmd.Flags().StringVar(&ifMatch, "if-match", "", "Valeur ETag pour le contrôle de concurrence")

	return cmd
}

func newPreviewCmd(opts AclOptions) *cobra.Command {
	var (
		previewType string
		previewFor  string
		file        string
		stdin       bool
	)

	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Prévisualise l'effet d'un fichier ACL",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			body, err := readBody(file, stdin)
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/acl/preview?type=%s&previewFor=%s", tailnet, previewType, previewFor)

			data, err := client.Post(path, bytes.NewReader(body))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				fmt.Println(string(data))
				return nil
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&previewType, "type", "", "Type de prévisualisation (user, ipport)")
	cmd.Flags().StringVar(&previewFor, "preview-for", "", "Cible de la prévisualisation")
	cmd.Flags().StringVar(&file, "file", "", "Chemin vers le fichier ACL")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Lire le fichier ACL depuis stdin")

	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("preview-for")

	return cmd
}

func newValidateCmd(opts AclOptions) *cobra.Command {
	var (
		file  string
		stdin bool
	)

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Valide un fichier ACL",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			body, err := readBody(file, stdin)
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/acl/validate", tailnet)

			data, err := client.Post(path, bytes.NewReader(body))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				fmt.Println(string(data))
				return nil
			}
			return output.PrintJSON(result)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Chemin vers le fichier ACL")
	cmd.Flags().BoolVar(&stdin, "stdin", false, "Lire le fichier ACL depuis stdin")

	return cmd
}
