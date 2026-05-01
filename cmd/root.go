package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dimer47/tailscale-cli/cmd/acl"
	"github.com/dimer47/tailscale-cli/cmd/auth"
	"github.com/dimer47/tailscale-cli/cmd/contact"
	"github.com/dimer47/tailscale-cli/cmd/device"
	"github.com/dimer47/tailscale-cli/cmd/dns"
	"github.com/dimer47/tailscale-cli/cmd/invite"
	"github.com/dimer47/tailscale-cli/cmd/key"
	cmdlog "github.com/dimer47/tailscale-cli/cmd/log"
	mcp_cmd "github.com/dimer47/tailscale-cli/cmd/mcp"
	"github.com/dimer47/tailscale-cli/cmd/posture"
	"github.com/dimer47/tailscale-cli/cmd/service"
	"github.com/dimer47/tailscale-cli/cmd/settings"
	"github.com/dimer47/tailscale-cli/cmd/user"
	"github.com/dimer47/tailscale-cli/cmd/webhook"
	"github.com/dimer47/tailscale-cli/internal/api"
	"github.com/dimer47/tailscale-cli/internal/config"
	"github.com/dimer47/tailscale-cli/internal/keychain"
	"github.com/dimer47/tailscale-cli/internal/update"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "tailscale-cli",
	Short: "CLI pour l'API Tailscale v2",
	Long: `tailscale-cli est un outil en ligne de commande pour interagir avec l'API Tailscale v2.

Il permet de gérer les appareils, les ACL, les clés d'authentification, le DNS,
les utilisateurs, les webhooks et bien d'autres ressources de votre tailnet,
le tout depuis votre terminal.

Configurez votre token API avec 'tailscale-cli auth login' ou via la variable
d'environnement TSCLI_API_TOKEN.

Le token est stocké de façon sécurisée dans le trousseau système.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags
	rootCmd.PersistentFlags().StringP("api-token", "t", "", "Token API Tailscale (env: TSCLI_API_TOKEN)")
	rootCmd.PersistentFlags().StringP("tailnet", "n", "-", "Tailnet cible (env: TSCLI_TAILNET)")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "Format de sortie : table, json, yaml (env: TSCLI_OUTPUT)")
	rootCmd.PersistentFlags().Bool("json", false, "Raccourci pour --output json")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Mode silencieux")
	rootCmd.PersistentFlags().Bool("debug", false, "Mode debug (env: TSCLI_DEBUG)")
	rootCmd.PersistentFlags().Bool("no-color", false, "Désactive les couleurs (env: NO_COLOR)")
	rootCmd.PersistentFlags().String("config", config.DefaultConfigPath(), "Chemin du fichier de configuration (env: TSCLI_CONFIG)")
	rootCmd.PersistentFlags().StringP("context", "c", "default", "Contexte de configuration à utiliser (env: TSCLI_CONTEXT)")

	// Viper bindings
	_ = viper.BindPFlag("api-token", rootCmd.PersistentFlags().Lookup("api-token"))
	_ = viper.BindPFlag("tailnet", rootCmd.PersistentFlags().Lookup("tailnet"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))

	// Environment variable bindings
	_ = viper.BindEnv("api-token", "TSCLI_API_TOKEN")
	_ = viper.BindEnv("tailnet", "TSCLI_TAILNET")
	_ = viper.BindEnv("output", "TSCLI_OUTPUT")
	_ = viper.BindEnv("debug", "TSCLI_DEBUG")
	_ = viper.BindEnv("no-color", "NO_COLOR")
	_ = viper.BindEnv("config", "TSCLI_CONFIG")
	_ = viper.BindEnv("context", "TSCLI_CONTEXT")

	// Sub-commands
	rootCmd.AddCommand(device.NewCmdDevice(&device.DeviceOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(acl.NewCmdAcl(acl.AclOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(key.NewCmdKey(key.KeyOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(dns.NewCmdDns(dns.DnsOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(cmdlog.NewCmdLog(cmdlog.LogOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(user.NewCmdUser(user.UserOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(invite.NewCmdInvite(invite.InviteOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(posture.NewCmdPosture(posture.PostureOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(contact.NewCmdContact(contact.ContactOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(webhook.NewCmdWebhook(webhook.WebhookOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(settings.NewCmdSettings(settings.SettingsOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(service.NewCmdService(service.ServiceOptions{
		GetClient:       getClient,
		GetOutputFormat: getOutputFormat,
		GetTailnet:      getTailnet,
	}))
	rootCmd.AddCommand(auth.NewCmdAuth())
	rootCmd.AddCommand(mcp_cmd.NewCmdMcpServe())
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newSelfUpdateCmd())
	rootCmd.AddCommand(newCompletionCmd())
}

func initConfig() {
	viper.SetEnvPrefix("TSCLI")
	viper.AutomaticEnv()
}

// Execute runs the root command.
func Execute() {
	// Launch update check in background (non-blocking)
	updateCh := make(chan *update.CheckResult, 1)
	go func() {
		updateCh <- update.Check(Version)
	}()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	// Print update notice if available (wait max 1s for the goroutine)
	select {
	case result := <-updateCh:
		if result != nil && result.UpdateAvailable {
			fmt.Fprint(os.Stderr, result.Message)
		}
	case <-time.After(1 * time.Second):
	}
}

// getClient resolves the API token and creates an API client.
// Priority: 1) --api-token flag  2) TSCLI_API_TOKEN env  3) system keyring  4) config file (legacy)
func getClient() (*api.Client, error) {
	token := viper.GetString("api-token")

	if token == "" {
		token = resolveTokenFromKeychain()
	}

	// Legacy fallback: token in config file (for migration)
	if token == "" {
		token = resolveTokenFromConfigFile()
	}

	if token == "" {
		return nil, fmt.Errorf("aucun token API configuré. Utilisez --api-token, TSCLI_API_TOKEN, ou 'tailscale-cli auth login'")
	}

	opts := []api.Option{
		api.WithDebug(viper.GetBool("debug")),
	}

	return api.NewClient(token, opts...), nil
}

// resolveTokenFromKeychain tries to get the token from the system keyring.
func resolveTokenFromKeychain() string {
	if !keychain.IsAvailable() {
		return ""
	}

	contextName := resolveContextName()
	token, err := keychain.Get(contextName)
	if err != nil {
		return ""
	}
	return token
}

// resolveTokenFromConfigFile tries to get the token from the config file (legacy).
func resolveTokenFromConfigFile() string {
	cfgPath := viper.GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return ""
	}

	contextName := viper.GetString("context")
	ctx, _, err := config.GetActiveContext(cfg, contextName)
	if err != nil {
		return ""
	}

	return ctx.APIToken
}

// resolveContextName resolves the effective context name from config.
func resolveContextName() string {
	contextName := viper.GetString("context")

	cfgPath := viper.GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		if contextName == "" {
			return "default"
		}
		return contextName
	}

	_, resolved, _ := config.GetActiveContext(cfg, contextName)
	return resolved
}

// getOutputFormat resolves the output format. --json overrides --output.
func getOutputFormat() string {
	if viper.GetBool("json") {
		return "json"
	}
	return viper.GetString("output")
}

// getTailnet resolves the tailnet. Falls back to config if not set via flag/env.
func getTailnet() string {
	tailnet := viper.GetString("tailnet")

	// Si le tailnet est la valeur par défaut, essayer la config
	if tailnet == "-" {
		cfgPath := viper.GetString("config")
		cfg, err := config.Load(cfgPath)
		if err == nil {
			contextName := viper.GetString("context")
			ctx, _, err := config.GetActiveContext(cfg, contextName)
			if err == nil && ctx.Tailnet != "" {
				return ctx.Tailnet
			}
		}
	}

	return tailnet
}
