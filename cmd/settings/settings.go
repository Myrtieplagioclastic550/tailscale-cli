package settings

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/output"
	"github.com/spf13/cobra"
)

// SettingsOptions contains the dependencies for the settings commands.
type SettingsOptions struct {
	GetClient       func() (*api.Client, error)
	GetOutputFormat func() string
	GetTailnet      func() string
}

// NewCmdSettings returns the settings command group.
func NewCmdSettings(opts SettingsOptions) *cobra.Command {
	settingsCmd := &cobra.Command{
		Use:   "settings",
		Short: "Gestion des paramètres du tailnet",
		Long:  "Commandes pour consulter et modifier les paramètres du tailnet.",
	}

	settingsCmd.AddCommand(newGetCmd(opts))
	settingsCmd.AddCommand(newUpdateCmd(opts))

	return settingsCmd
}

func newGetCmd(opts SettingsOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Affiche les paramètres du tailnet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			path := fmt.Sprintf("/tailnet/%s/settings", opts.GetTailnet())
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

func newUpdateCmd(opts SettingsOptions) *cobra.Command {
	var devicesApproval bool
	var devicesAutoUpdates bool
	var devicesKeyDuration int
	var usersApproval bool
	var usersExternalTailnets string
	var networkFlowLogging bool
	var regionalRouting bool
	var postureIdentityCollection bool
	var httpsEnabled bool
	var aclsExternal bool
	var aclsExternalLink string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Met à jour les paramètres du tailnet",
		Long:  "Met à jour les paramètres du tailnet. Seuls les flags explicitement fournis sont envoyés.",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := opts.GetClient()
			if err != nil {
				return err
			}

			payload := map[string]interface{}{}

			// Devices
			devices := map[string]interface{}{}
			if cmd.Flags().Changed("devices-approval") {
				devices["approvalRequired"] = devicesApproval
			}
			if cmd.Flags().Changed("devices-auto-updates") {
				devices["autoUpdatesEnabled"] = devicesAutoUpdates
			}
			if cmd.Flags().Changed("devices-key-duration") {
				devices["keyDurationDays"] = devicesKeyDuration
			}
			if len(devices) > 0 {
				payload["devices"] = devices
			}

			// Users
			users := map[string]interface{}{}
			if cmd.Flags().Changed("users-approval") {
				users["approvalRequired"] = usersApproval
			}
			if cmd.Flags().Changed("users-external-tailnets") {
				users["externalTailnets"] = usersExternalTailnets
			}
			if len(users) > 0 {
				payload["users"] = users
			}

			// Network
			network := map[string]interface{}{}
			if cmd.Flags().Changed("network-flow-logging") {
				network["flowLoggingEnabled"] = networkFlowLogging
			}
			if cmd.Flags().Changed("regional-routing") {
				network["regionalRoutingEnabled"] = regionalRouting
			}
			if len(network) > 0 {
				payload["network"] = network
			}

			// Posture
			if cmd.Flags().Changed("posture-identity-collection") {
				payload["postureIdentityCollection"] = postureIdentityCollection
			}

			// HTTPS
			if cmd.Flags().Changed("https") {
				payload["httpsEnabled"] = httpsEnabled
			}

			// ACLs
			acls := map[string]interface{}{}
			if cmd.Flags().Changed("acls-external") {
				acls["externallyManaged"] = aclsExternal
			}
			if cmd.Flags().Changed("acls-external-link") {
				acls["externalLink"] = aclsExternalLink
			}
			if len(acls) > 0 {
				payload["acls"] = acls
			}

			if len(payload) == 0 {
				return fmt.Errorf("aucun paramètre à mettre à jour, utilisez au moins un flag")
			}

			data, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}

			path := fmt.Sprintf("/tailnet/%s/settings", opts.GetTailnet())
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

	cmd.Flags().BoolVar(&devicesApproval, "devices-approval", false, "Activer/désactiver l'approbation des appareils")
	cmd.Flags().BoolVar(&devicesAutoUpdates, "devices-auto-updates", false, "Activer/désactiver les mises à jour automatiques")
	cmd.Flags().IntVar(&devicesKeyDuration, "devices-key-duration", 0, "Durée de validité des clés (en jours)")
	cmd.Flags().BoolVar(&usersApproval, "users-approval", false, "Activer/désactiver l'approbation des utilisateurs")
	cmd.Flags().StringVar(&usersExternalTailnets, "users-external-tailnets", "", "Politique des tailnets externes (none, admin, member)")
	cmd.Flags().BoolVar(&networkFlowLogging, "network-flow-logging", false, "Activer/désactiver la journalisation des flux réseau")
	cmd.Flags().BoolVar(&regionalRouting, "regional-routing", false, "Activer/désactiver le routage régional")
	cmd.Flags().BoolVar(&postureIdentityCollection, "posture-identity-collection", false, "Activer/désactiver la collecte d'identité posture")
	cmd.Flags().BoolVar(&httpsEnabled, "https", false, "Activer/désactiver HTTPS")
	cmd.Flags().BoolVar(&aclsExternal, "acls-external", false, "Activer/désactiver la gestion externe des ACLs")
	cmd.Flags().StringVar(&aclsExternalLink, "acls-external-link", "", "Lien vers les ACLs externes")

	return cmd
}
