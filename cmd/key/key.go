package key

import (
	"encoding/json"
	"fmt"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// KeyOptions contains closures for lazy resolution of shared options.
type KeyOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdKey returns the key parent command with all its subcommands.
func NewCmdKey(opts KeyOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Gestion des clés d'authentification",
		Long:  "Commandes pour créer, lister, consulter, modifier et supprimer les clés d'authentification du tailnet.",
	}

	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newCreateCmd(opts))
	cmd.AddCommand(newGetCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))

	return cmd
}

func newListCmd(opts KeyOptions) *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Liste les clés du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			path := fmt.Sprintf("/tailnet/%s/keys", tailnet)
			if all {
				path += "?all=true"
			}

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

	cmd.Flags().BoolVar(&all, "all", false, "Inclure toutes les clés (pas uniquement les siennes)")

	return cmd
}

func newCreateCmd(opts KeyOptions) *cobra.Command {
	var (
		keyType       string
		description   string
		expiry        int64
		reusable      bool
		ephemeral     bool
		preauthorized bool
		tags          []string
		scopes        []string
		issuer        string
		subject       string
		audience      string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Crée une nouvelle clé",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()

			body := map[string]interface{}{
				"keyType": keyType,
			}

			if description != "" {
				body["description"] = description
			}

			switch keyType {
			case "auth":
				capabilities := map[string]interface{}{}
				devices := map[string]interface{}{
					"create": map[string]interface{}{
						"reusable":      reusable,
						"ephemeral":     ephemeral,
						"preauthorized": preauthorized,
						"tags":          tags,
					},
				}
				capabilities["devices"] = devices
				body["capabilities"] = capabilities
				if expiry > 0 {
					body["expirySeconds"] = expiry
				}
			case "client":
				if len(scopes) > 0 {
					body["scopes"] = scopes
				}
				if expiry > 0 {
					body["expirySeconds"] = expiry
				}
			case "federated":
				federated := map[string]interface{}{}
				if issuer != "" {
					federated["issuer"] = issuer
				}
				if subject != "" {
					federated["subject"] = subject
				}
				if audience != "" {
					federated["audience"] = audience
				}
				body["federated"] = federated
				if len(scopes) > 0 {
					body["scopes"] = scopes
				}
			}

			path := fmt.Sprintf("/tailnet/%s/keys", tailnet)

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

	cmd.Flags().StringVar(&keyType, "type", "auth", "Type de clé (auth, client, federated)")
	cmd.Flags().StringVar(&description, "description", "", "Description de la clé")
	cmd.Flags().Int64Var(&expiry, "expiry", 0, "Durée d'expiration en secondes")
	cmd.Flags().BoolVar(&reusable, "reusable", false, "Clé réutilisable (type auth)")
	cmd.Flags().BoolVar(&ephemeral, "ephemeral", false, "Clé éphémère (type auth)")
	cmd.Flags().BoolVar(&preauthorized, "preauthorized", false, "Clé pré-autorisée (type auth)")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags associés à la clé (type auth)")
	cmd.Flags().StringSliceVar(&scopes, "scopes", nil, "Scopes de la clé (type client/federated)")
	cmd.Flags().StringVar(&issuer, "issuer", "", "Émetteur (type federated)")
	cmd.Flags().StringVar(&subject, "subject", "", "Sujet (type federated)")
	cmd.Flags().StringVar(&audience, "audience", "", "Audience (type federated)")

	return cmd
}

func newGetCmd(opts KeyOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [keyId]",
		Short: "Récupère les détails d'une clé",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()
			keyID := args[0]

			path := fmt.Sprintf("/tailnet/%s/keys/%s", tailnet, keyID)

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

	return cmd
}

func newDeleteCmd(opts KeyOptions) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete [keyId]",
		Short: "Supprime une clé",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("utilisez --confirm pour confirmer la suppression")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()
			keyID := args[0]

			path := fmt.Sprintf("/tailnet/%s/keys/%s", tailnet, keyID)

			_, err = client.Delete(path)
			if err != nil {
				return err
			}

			fmt.Printf("Clé %s supprimée avec succès.\n", keyID)
			return nil
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirmer la suppression")

	return cmd
}

func newUpdateCmd(opts KeyOptions) *cobra.Command {
	var (
		description string
		scopes      []string
		tags        []string
		issuer      string
		subject     string
		audience    string
	)

	cmd := &cobra.Command{
		Use:   "update [keyId]",
		Short: "Met à jour une clé existante",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}
			tailnet := opts.GetTailnet()
			keyID := args[0]

			body := map[string]interface{}{}

			if cmd.Flags().Changed("description") {
				body["description"] = description
			}
			if cmd.Flags().Changed("scopes") {
				body["scopes"] = scopes
			}
			if cmd.Flags().Changed("tags") {
				body["tags"] = tags
			}
			if cmd.Flags().Changed("issuer") {
				body["issuer"] = issuer
			}
			if cmd.Flags().Changed("subject") {
				body["subject"] = subject
			}
			if cmd.Flags().Changed("audience") {
				body["audience"] = audience
			}

			path := fmt.Sprintf("/tailnet/%s/keys/%s", tailnet, keyID)

			data, err := client.DoJSON("PUT", path, body)
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

	cmd.Flags().StringVar(&description, "description", "", "Description de la clé")
	cmd.Flags().StringSliceVar(&scopes, "scopes", nil, "Scopes de la clé")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags associés à la clé")
	cmd.Flags().StringVar(&issuer, "issuer", "", "Émetteur")
	cmd.Flags().StringVar(&subject, "subject", "", "Sujet")
	cmd.Flags().StringVar(&audience, "audience", "", "Audience")

	return cmd
}
