package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// ServiceOptions holds the dependencies for service commands.
type ServiceOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdService returns the service command group with all subcommands.
func NewCmdService(opts ServiceOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Gestion des services VIP",
		Long:  "Commandes pour lister, inspecter, creer, modifier et supprimer les services VIP du tailnet.",
	}

	cmd.AddCommand(newCmdServiceList(opts))
	cmd.AddCommand(newCmdServiceGet(opts))
	cmd.AddCommand(newCmdServiceCreate(opts))
	cmd.AddCommand(newCmdServiceUpdate(opts))
	cmd.AddCommand(newCmdServiceDelete(opts))
	cmd.AddCommand(newCmdServiceHosts(opts))
	cmd.AddCommand(newCmdServiceApprove(opts))

	return cmd
}

// ---------------------------------------------------------------------------
// list
// ---------------------------------------------------------------------------

func newCmdServiceList(opts ServiceOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lister tous les services VIP du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			tailnet := opts.GetTailnet()
			path := fmt.Sprintf("/tailnet/%s/services", tailnet)

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

// ---------------------------------------------------------------------------
// get
// ---------------------------------------------------------------------------

func newCmdServiceGet(opts ServiceOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <serviceName>",
		Short: "Afficher les details d'un service VIP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			tailnet := opts.GetTailnet()
			path := fmt.Sprintf("/tailnet/%s/services/%s", tailnet, serviceName)

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

// ---------------------------------------------------------------------------
// create
// ---------------------------------------------------------------------------

func newCmdServiceCreate(opts ServiceOptions) *cobra.Command {
	var (
		comment string
		ports   []string
		tags    []string
		ipv4    string
	)

	cmd := &cobra.Command{
		Use:   "create <serviceName>",
		Short: "Creer un service VIP",
		Long:  "Creer un service VIP. Le nom doit commencer par \"svc:\".",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			if !strings.HasPrefix(serviceName, "svc:") {
				return fmt.Errorf("le nom du service doit commencer par \"svc:\", recu: %s", serviceName)
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			tailnet := opts.GetTailnet()
			path := fmt.Sprintf("/tailnet/%s/services/%s", tailnet, serviceName)

			payload := map[string]interface{}{
				"name":    serviceName,
				"comment": comment,
				"ports":   ports,
				"tags":    tags,
			}
			if ipv4 != "" {
				payload["addrs"] = []string{ipv4}
			}

			jsonBody, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			body, err := client.Put(path, bytes.NewReader(jsonBody))
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

	cmd.Flags().StringVar(&comment, "comment", "", "Commentaire du service")
	cmd.Flags().StringSliceVar(&ports, "ports", nil, "Ports exposes par le service")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags associes au service")
	cmd.Flags().StringVar(&ipv4, "ipv4", "", "Adresse IPv4 du service")

	return cmd
}

// ---------------------------------------------------------------------------
// update
// ---------------------------------------------------------------------------

func newCmdServiceUpdate(opts ServiceOptions) *cobra.Command {
	var (
		comment string
		ports   []string
		tags    []string
		ipv4    string
		newName string
	)

	cmd := &cobra.Command{
		Use:   "update <serviceName>",
		Short: "Mettre a jour un service VIP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			tailnet := opts.GetTailnet()
			path := fmt.Sprintf("/tailnet/%s/services/%s", tailnet, serviceName)

			payload := map[string]interface{}{}

			if cmd.Flags().Changed("new-name") {
				payload["name"] = newName
			} else {
				payload["name"] = serviceName
			}

			if cmd.Flags().Changed("comment") {
				payload["comment"] = comment
			}
			if cmd.Flags().Changed("ports") {
				payload["ports"] = ports
			}
			if cmd.Flags().Changed("tags") {
				payload["tags"] = tags
			}
			if cmd.Flags().Changed("ipv4") {
				payload["addrs"] = []string{ipv4}
			}

			jsonBody, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			body, err := client.Put(path, bytes.NewReader(jsonBody))
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

	cmd.Flags().StringVar(&comment, "comment", "", "Commentaire du service")
	cmd.Flags().StringSliceVar(&ports, "ports", nil, "Ports exposes par le service")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Tags associes au service")
	cmd.Flags().StringVar(&ipv4, "ipv4", "", "Adresse IPv4 du service")
	cmd.Flags().StringVar(&newName, "new-name", "", "Nouveau nom du service (renommage)")

	return cmd
}

// ---------------------------------------------------------------------------
// delete
// ---------------------------------------------------------------------------

func newCmdServiceDelete(opts ServiceOptions) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <serviceName>",
		Short: "Supprimer un service VIP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			if !confirm {
				fmt.Fprintf(os.Stderr, "Etes-vous sur de vouloir supprimer le service %q ? (oui/non) : ", serviceName)
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
					if answer != "oui" && answer != "o" && answer != "yes" && answer != "y" {
						fmt.Fprintln(os.Stderr, "Suppression annulee.")
						return nil
					}
				}
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			tailnet := opts.GetTailnet()
			path := fmt.Sprintf("/tailnet/%s/services/%s", tailnet, serviceName)

			_, err = client.Delete(path)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Service %q supprime avec succes.\n", serviceName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirmer la suppression sans prompt interactif")

	return cmd
}

// ---------------------------------------------------------------------------
// hosts
// ---------------------------------------------------------------------------

func newCmdServiceHosts(opts ServiceOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "hosts <serviceName>",
		Short: "Lister les appareils hebergeant un service VIP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			tailnet := opts.GetTailnet()
			path := fmt.Sprintf("/tailnet/%s/services/%s/devices", tailnet, serviceName)

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

// ---------------------------------------------------------------------------
// approve
// ---------------------------------------------------------------------------

func newCmdServiceApprove(opts ServiceOptions) *cobra.Command {
	var approved bool

	cmd := &cobra.Command{
		Use:   "approve <serviceName> <deviceId>",
		Short: "Approuver ou rejeter un appareil pour un service VIP",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			deviceID := args[1]

			if !cmd.Flags().Changed("approved") {
				return fmt.Errorf("le flag --approved est requis")
			}

			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			tailnet := opts.GetTailnet()
			path := fmt.Sprintf("/tailnet/%s/services/%s/device/%s/approved", tailnet, serviceName, deviceID)

			payload := map[string]interface{}{
				"approved": approved,
			}

			jsonBody, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			body, err := client.Post(path, bytes.NewReader(jsonBody))
			if err != nil {
				return err
			}

			if len(body) > 0 {
				var data interface{}
				if err := json.Unmarshal(body, &data); err != nil {
					return fmt.Errorf("parsing response: %w", err)
				}
				return output.Print(opts.GetOutputFormat(), data, nil)
			}

			status := "approuve"
			if !approved {
				status = "rejete"
			}
			fmt.Fprintf(os.Stderr, "Appareil %q %s pour le service %q.\n", deviceID, status, serviceName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&approved, "approved", false, "Approuver (true) ou rejeter (false) l'appareil")

	return cmd
}
