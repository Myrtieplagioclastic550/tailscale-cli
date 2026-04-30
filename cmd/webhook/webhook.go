package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// WebhookOptions contains the dependencies for the webhook commands.
type WebhookOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdWebhook returns the webhook command group.
func NewCmdWebhook(opts WebhookOptions) *cobra.Command {
	webhookCmd := &cobra.Command{
		Use:   "webhook",
		Short: "Gestion des webhooks",
		Long:  "Commandes pour créer, lister et gérer les webhooks du tailnet.",
	}

	webhookCmd.AddCommand(newListCmd(opts))
	webhookCmd.AddCommand(newCreateCmd(opts))
	webhookCmd.AddCommand(newGetCmd(opts))
	webhookCmd.AddCommand(newUpdateCmd(opts))
	webhookCmd.AddCommand(newDeleteCmd(opts))
	webhookCmd.AddCommand(newTestCmd(opts))
	webhookCmd.AddCommand(newRotateSecretCmd(opts))

	return webhookCmd
}

func newListCmd(opts WebhookOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Liste les webhooks du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/webhooks", opts.GetTailnet())
			body, err := client.Get(path)
			if err != nil {
				return err
			}

			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), data, nil)
		},
	}
}

func newCreateCmd(opts WebhookOptions) *cobra.Command {
	var url string
	var provider string
	var events []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Crée un webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			if url == "" {
				return fmt.Errorf("le flag --url est requis")
			}
			if len(events) == 0 {
				return fmt.Errorf("le flag --events est requis")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]interface{}{
				"endpointUrl":   url,
				"subscriptions": events,
			}
			if provider != "" {
				payload["providerType"] = provider
			}

			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/tailnet/%s/webhooks", opts.GetTailnet())
			body, err := client.Post(path, bytes.NewReader(data))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), result, nil)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "URL du webhook (requis)")
	cmd.Flags().StringVar(&provider, "provider", "", "Type de fournisseur (slack, mattermost, googlechat, discord)")
	cmd.Flags().StringSliceVar(&events, "events", nil, "Événements à écouter (requis)")

	return cmd
}

func newGetCmd(opts WebhookOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Affiche les détails d'un webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/webhooks/%s", args[0])
			body, err := client.Get(path)
			if err != nil {
				return err
			}

			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), data, nil)
		},
	}
}

func newUpdateCmd(opts WebhookOptions) *cobra.Command {
	var events []string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Met à jour un webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]interface{}{}
			if cmd.Flags().Changed("events") {
				payload["subscriptions"] = events
			}

			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/webhooks/%s", args[0])
			body, err := client.Patch(path, bytes.NewReader(data))
			if err != nil {
				return err
			}

			var result interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				return fmt.Errorf("parsing response: %w", err)
			}

			return output.Print(opts.GetOutputFormat(), result, nil)
		},
	}

	cmd.Flags().StringSliceVar(&events, "events", nil, "Événements à écouter")

	return cmd
}

func newDeleteCmd(opts WebhookOptions) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Supprime un webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("veuillez confirmer la suppression avec --confirm")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/webhooks/%s", args[0])
			body, err := client.Delete(path)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					return fmt.Errorf("parsing response: %w", err)
				}
				return output.Print(opts.GetOutputFormat(), result, nil)
			}

			fmt.Fprintf(os.Stderr, "Webhook %s supprimé.\n", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirmer la suppression")

	return cmd
}

func newTestCmd(opts WebhookOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "test <id>",
		Short: "Teste un webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/webhooks/%s/test", args[0])
			body, err := client.Post(path, nil)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					return fmt.Errorf("parsing response: %w", err)
				}
				return output.Print(opts.GetOutputFormat(), result, nil)
			}

			fmt.Fprintf(os.Stderr, "Test du webhook %s envoyé.\n", args[0])
			return nil
		},
	}
}

func newRotateSecretCmd(opts WebhookOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "rotate-secret <id>",
		Short: "Effectue la rotation du secret d'un webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/webhooks/%s/rotate", args[0])
			body, err := client.Post(path, nil)
			if err != nil {
				return err
			}

			if len(body) > 0 {
				var result interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					return fmt.Errorf("parsing response: %w", err)
				}
				return output.Print(opts.GetOutputFormat(), result, nil)
			}

			fmt.Fprintf(os.Stderr, "Secret du webhook %s renouvelé.\n", args[0])
			return nil
		},
	}
}
