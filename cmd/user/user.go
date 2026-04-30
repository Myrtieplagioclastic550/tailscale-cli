package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// UserOptions contains the dependencies for the user commands.
type UserOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdUser returns the user command group.
func NewCmdUser(opts UserOptions) *cobra.Command {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Gestion des utilisateurs",
		Long:  "Commandes pour lister, inspecter et gérer les utilisateurs du tailnet.",
	}

	userCmd.AddCommand(newListCmd(opts))
	userCmd.AddCommand(newGetCmd(opts))
	userCmd.AddCommand(newSetRoleCmd(opts))
	userCmd.AddCommand(newApproveCmd(opts))
	userCmd.AddCommand(newSuspendCmd(opts))
	userCmd.AddCommand(newRestoreCmd(opts))
	userCmd.AddCommand(newDeleteCmd(opts))

	return userCmd
}

func newListCmd(opts UserOptions) *cobra.Command {
	var userType string
	var role string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Liste les utilisateurs du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/users", opts.GetTailnet())

			sep := "?"
			if userType != "" {
				path += fmt.Sprintf("%stype=%s", sep, userType)
				sep = "&"
			}
			if role != "" {
				path += fmt.Sprintf("%srole=%s", sep, role)
			}

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

	cmd.Flags().StringVar(&userType, "type", "", "Type d'utilisateur (member, shared, all)")
	cmd.Flags().StringVar(&role, "role", "", "Filtrer par rôle")

	return cmd
}

func newGetCmd(opts UserOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <userId>",
		Short: "Affiche les détails d'un utilisateur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/users/%s", args[0])
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

func newSetRoleCmd(opts UserOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "set-role <userId> <role>",
		Short: "Change le rôle d'un utilisateur",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]string{"role": args[1]}
			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/users/%s/role", args[0])
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
}

func newApproveCmd(opts UserOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "approve <userId>",
		Short: "Approuve un utilisateur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/users/%s/approve", args[0])
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

			fmt.Fprintf(os.Stderr, "Utilisateur %s approuvé.\n", args[0])
			return nil
		},
	}
}

func newSuspendCmd(opts UserOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "suspend <userId>",
		Short: "Suspend un utilisateur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/users/%s/suspend", args[0])
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

			fmt.Fprintf(os.Stderr, "Utilisateur %s suspendu.\n", args[0])
			return nil
		},
	}
}

func newRestoreCmd(opts UserOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "restore <userId>",
		Short: "Restaure un utilisateur suspendu",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/users/%s/restore", args[0])
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

			fmt.Fprintf(os.Stderr, "Utilisateur %s restauré.\n", args[0])
			return nil
		},
	}
}

func newDeleteCmd(opts UserOptions) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <userId>",
		Short: "Supprime un utilisateur",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("veuillez confirmer la suppression avec --confirm")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/users/%s/delete", args[0])
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

			fmt.Fprintf(os.Stderr, "Utilisateur %s supprimé.\n", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirmer la suppression")

	return cmd
}
