package contact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// ContactOptions contains the dependencies for the contact commands.
type ContactOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdContact returns the contact command group.
func NewCmdContact(opts ContactOptions) *cobra.Command {
	contactCmd := &cobra.Command{
		Use:   "contact",
		Short: "Gestion des contacts du tailnet",
		Long:  "Commandes pour consulter et modifier les contacts du tailnet.",
	}

	contactCmd.AddCommand(newGetCmd(opts))
	contactCmd.AddCommand(newUpdateCmd(opts))
	contactCmd.AddCommand(newResendVerificationCmd(opts))

	return contactCmd
}

func newGetCmd(opts ContactOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Affiche les contacts du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/contacts", opts.GetTailnet())
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

func newUpdateCmd(opts ContactOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "update <type> <email>",
		Short: "Met à jour un contact du tailnet",
		Long:  "Met à jour un contact. Le type doit être account, support ou security.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			contactType := args[0]
			email := args[1]

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]string{"email": email}
			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/tailnet/%s/contacts/%s", opts.GetTailnet(), contactType)
			body, err := client.Patch(path, bytes.NewReader(data))
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

			fmt.Fprintf(os.Stderr, "Contact %s mis à jour avec l'email %s.\n", contactType, email)
			return nil
		},
	}
}

func newResendVerificationCmd(opts ContactOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "resend-verification <type>",
		Short: "Renvoie l'email de vérification pour un contact",
		Long:  "Renvoie l'email de vérification. Le type doit être account, support ou security.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contactType := args[0]

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/contacts/%s/resend-verification-email", opts.GetTailnet(), contactType)
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

			fmt.Fprintf(os.Stderr, "Email de vérification renvoyé pour le contact %s.\n", contactType)
			return nil
		},
	}
}
