package posture

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// PostureOptions contains the dependencies for the posture commands.
type PostureOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdPosture returns the posture command group.
func NewCmdPosture(opts PostureOptions) *cobra.Command {
	postureCmd := &cobra.Command{
		Use:   "posture",
		Short: "Gestion des intégrations de posture",
		Long:  "Commandes pour gérer les intégrations de posture du tailnet.",
	}

	postureCmd.AddCommand(newListCmd(opts))
	postureCmd.AddCommand(newCreateCmd(opts))
	postureCmd.AddCommand(newGetCmd(opts))
	postureCmd.AddCommand(newUpdateCmd(opts))
	postureCmd.AddCommand(newDeleteCmd(opts))

	return postureCmd
}

func newListCmd(opts PostureOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Liste les intégrations de posture",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/posture/integrations", opts.GetTailnet())
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

func newCreateCmd(opts PostureOptions) *cobra.Command {
	var provider string
	var cloudID string
	var clientID string
	var tenantID string
	var clientSecret string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Crée une intégration de posture",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]string{
				"provider":     provider,
				"cloudId":      cloudID,
				"clientId":     clientID,
				"tenantId":     tenantID,
				"clientSecret": clientSecret,
			}
			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/tailnet/%s/posture/integrations", opts.GetTailnet())
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

	cmd.Flags().StringVar(&provider, "provider", "", "Fournisseur de posture")
	cmd.Flags().StringVar(&cloudID, "cloud-id", "", "Identifiant cloud")
	cmd.Flags().StringVar(&clientID, "client-id", "", "Identifiant client")
	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Identifiant du tenant")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "Secret client")

	return cmd
}

func newGetCmd(opts PostureOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Affiche les détails d'une intégration de posture",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/posture/integrations/%s", args[0])
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

func newUpdateCmd(opts PostureOptions) *cobra.Command {
	var cloudID string
	var clientID string
	var tenantID string
	var clientSecret string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Met à jour une intégration de posture",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]string{}
			if cmd.Flags().Changed("cloud-id") {
				payload["cloudId"] = cloudID
			}
			if cmd.Flags().Changed("client-id") {
				payload["clientId"] = clientID
			}
			if cmd.Flags().Changed("tenant-id") {
				payload["tenantId"] = tenantID
			}
			if cmd.Flags().Changed("client-secret") {
				payload["clientSecret"] = clientSecret
			}

			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/posture/integrations/%s", args[0])
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

	cmd.Flags().StringVar(&cloudID, "cloud-id", "", "Identifiant cloud")
	cmd.Flags().StringVar(&clientID, "client-id", "", "Identifiant client")
	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Identifiant du tenant")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "Secret client")

	return cmd
}

func newDeleteCmd(opts PostureOptions) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Supprime une intégration de posture",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("veuillez confirmer la suppression avec --confirm")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/posture/integrations/%s", args[0])
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

			fmt.Fprintf(os.Stderr, "Intégration de posture %s supprimée.\n", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirmer la suppression")

	return cmd
}
