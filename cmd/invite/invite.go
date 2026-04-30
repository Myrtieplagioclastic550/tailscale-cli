package invite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// InviteOptions contains the dependencies for the invite commands.
type InviteOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdInvite returns the invite command group.
func NewCmdInvite(opts InviteOptions) *cobra.Command {
	inviteCmd := &cobra.Command{
		Use:   "invite",
		Short: "Gestion des invitations",
		Long:  "Commandes pour créer, lister et gérer les invitations au tailnet.",
	}

	inviteCmd.AddCommand(newUserCmd(opts))
	inviteCmd.AddCommand(newDeviceCmd(opts))

	return inviteCmd
}

// --- invite user ---

func newUserCmd(opts InviteOptions) *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Gestion des invitations utilisateur",
	}

	userCmd.AddCommand(newUserListCmd(opts))
	userCmd.AddCommand(newUserCreateCmd(opts))
	userCmd.AddCommand(newUserGetCmd(opts))
	userCmd.AddCommand(newUserDeleteCmd(opts))
	userCmd.AddCommand(newUserResendCmd(opts))

	return userCmd
}

func newUserListCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Liste les invitations utilisateur du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/user-invites", opts.GetTailnet())
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

func newUserCreateCmd(opts InviteOptions) *cobra.Command {
	var email string
	var role string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Crée une invitation utilisateur",
		RunE: func(cmd *cobra.Command, args []string) error {
			if email == "" {
				return fmt.Errorf("le flag --email est requis")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := []map[string]string{
				{
					"email": email,
					"role":  role,
				},
			}
			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/tailnet/%s/user-invites", opts.GetTailnet())
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

	cmd.Flags().StringVar(&email, "email", "", "Adresse email de l'invité")
	cmd.Flags().StringVar(&role, "role", "member", "Rôle de l'invité (default: member)")

	return cmd
}

func newUserGetCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Affiche les détails d'une invitation utilisateur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/user-invites/%s", args[0])
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

func newUserDeleteCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Supprime une invitation utilisateur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/user-invites/%s", args[0])
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

			fmt.Fprintf(os.Stderr, "Invitation utilisateur %s supprimée.\n", args[0])
			return nil
		},
	}
}

func newUserResendCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "resend <id>",
		Short: "Renvoie une invitation utilisateur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/user-invites/%s/resend", args[0])
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

			fmt.Fprintf(os.Stderr, "Invitation utilisateur %s renvoyée.\n", args[0])
			return nil
		},
	}
}

// --- invite device ---

func newDeviceCmd(opts InviteOptions) *cobra.Command {
	deviceCmd := &cobra.Command{
		Use:   "device",
		Short: "Gestion des invitations d'appareil",
	}

	deviceCmd.AddCommand(newDeviceListCmd(opts))
	deviceCmd.AddCommand(newDeviceCreateCmd(opts))
	deviceCmd.AddCommand(newDeviceGetCmd(opts))
	deviceCmd.AddCommand(newDeviceDeleteCmd(opts))
	deviceCmd.AddCommand(newDeviceResendCmd(opts))
	deviceCmd.AddCommand(newDeviceAcceptCmd(opts))

	return deviceCmd
}

func newDeviceListCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list <deviceId>",
		Short: "Liste les invitations d'un appareil",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/device/%s/device-invites", args[0])
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

func newDeviceCreateCmd(opts InviteOptions) *cobra.Command {
	var email string
	var multiUse bool
	var allowExitNode bool

	cmd := &cobra.Command{
		Use:   "create <deviceId>",
		Short: "Crée une invitation d'appareil",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			invite := map[string]interface{}{}
			if email != "" {
				invite["email"] = email
			}
			if multiUse {
				invite["multiUse"] = multiUse
			}
			if allowExitNode {
				invite["allowExitNode"] = allowExitNode
			}

			payload := []map[string]interface{}{invite}
			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/device/%s/device-invites", args[0])
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

	cmd.Flags().StringVar(&email, "email", "", "Adresse email de l'invité")
	cmd.Flags().BoolVar(&multiUse, "multi-use", false, "Invitation à usage multiple")
	cmd.Flags().BoolVar(&allowExitNode, "allow-exit-node", false, "Autoriser l'utilisation comme noeud de sortie")

	return cmd
}

func newDeviceGetCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Affiche les détails d'une invitation d'appareil",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/device-invites/%s", args[0])
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

func newDeviceDeleteCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Supprime une invitation d'appareil",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/device-invites/%s", args[0])
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

			fmt.Fprintf(os.Stderr, "Invitation d'appareil %s supprimée.\n", args[0])
			return nil
		},
	}
}

func newDeviceResendCmd(opts InviteOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "resend <id>",
		Short: "Renvoie une invitation d'appareil",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/device-invites/%s/resend", args[0])
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

			fmt.Fprintf(os.Stderr, "Invitation d'appareil %s renvoyée.\n", args[0])
			return nil
		},
	}
}

func newDeviceAcceptCmd(opts InviteOptions) *cobra.Command {
	var inviteURL string

	cmd := &cobra.Command{
		Use:   "accept",
		Short: "Accepte une invitation d'appareil",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inviteURL == "" {
				return fmt.Errorf("le flag --invite-url est requis")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]string{"invite": inviteURL}
			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := "/device-invites/-/accept"
			body, err := client.Post(path, bytes.NewReader(data))
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

			fmt.Fprintf(os.Stderr, "Invitation acceptée.\n")
			return nil
		},
	}

	cmd.Flags().StringVar(&inviteURL, "invite-url", "", "URL de l'invitation à accepter (requis)")

	return cmd
}
